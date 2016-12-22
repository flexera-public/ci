
# Not an executable. This is a function library that should be sourced
# by bash scripts. It contains a main() function which you can call from
# your own script.
#

# Automatically exit on error
set -e

# Before sourcing this script, you MUST define the following:
#    - a function before_build()

type before_build | grep -q 'before_build is a function'
if [ $? != 0 ]
then
  echo "ERROR: must define the before_build() function before sourcing docker-shared.sh"
  exit 20
fi

# Deduce some information about the Docker repository and build args
org_name=rightscale
app_name=`basename $PWD`
gitref=`git rev-parse --verify HEAD`

if [ $? != 0 ]
then
  echo "ERROR: working directory must be a git repository"
  exit 30
fi

# Login to DockerHub
login ()
{
  if [ -z "${DOCKERHUB_USER}${DOCKERHUB_PASSWORD}" ]
  then
    echo "ERROR: must export DOCKERHUB_(PASSWORD|USER) to perform this action"
    exit 40
  else
    echo "Logging into DockerHub as $DOCKERHUB_USER"
    docker login -u $DOCKERHUB_USER -p $DOCKERHUB_PASSWORD
  fi
}

# Build a docker image, including performing any prerequisite work for the build: gathering
# dependencies, creating intermediate build artifacts in a build container, etc.
build ()
{
  before_build

  echo "Building Docker image $org_name/$app_name:$1"
  docker build --build-arg gitref=$gitref --tag $org_name/$app_name:$1 .

  return $?
}

# Clean up intermediate build artifacts
clean ()
{
  echo "Removing intermediate build artifacts"
  rm -Rf vendor/cache vendor/bundle tmp/* log/*
}

# Push a named tag of this repo's image to DockerHub. This is a nearly-useless shortcut.
push ()
{
  echo "Pushing Docker image $org_name/$app_name:$1"
  docker push $org_name/$app_name:$1

  return $?
}

# Perform a continuous integration build. This involves the following steps:
#   - decide whether we need an image based on $1 (tag name) and $2 (pull request number)
#   - download bare Docker binary into bin; add bin to path
#   - build and push an image
ci ()
{
  if [[ "$2" != "false" ]]
  then
    echo "Skipping Docker image build due to pull-request status ($2)"
  elif [[ ! "$1" =~ (^(latest|staging|production-|experimental))|(_cow$) ]]
  then
    echo "Skipping Docker image build due to uninteresting branch name ($1)"
  else
    # Check if TRAVIS_COMMIT is already in Docker's image git.ref
    token=`curl -H 'Accept: application/json' --user "$DOCKERHUB_USER:$DOCKERHUB_PASSWORD" "https://auth.docker.io/token?service=registry.docker.io&scope=repository:${org_name}/${app_name}:pull" | jq --raw-output .token`
    image_git_ref=`curl -H "Authorization: Bearer $token" "https://registry-1.docker.io/v2/${org_name}/${app_name}/manifests/$1" | jq --raw-output .history[0].v1Compatibility | jq --raw-output '.config.Labels["git.ref"]'`
    if [[ $image_git_ref == $TRAVIS_COMMIT ]] ; then
      echo "Skipping Docker image build due to build's commit sha (${image_git_ref}) is equal to current $1 Docker image git.ref"
    else
      echo "Building $1 tag since current git.ref (${image_git_ref}) doesn't match build's commit sha (${TRAVIS_COMMIT})"
      login
      build $1 && push $1
    fi
  fi
}

# Print usage information
help ()
{
  echo
  echo "Usage:"
  echo " $0 login"
  echo " $0 build [tag]"
  echo " $0 clean"
  echo " $0 push [tag]"
  echo " $0 ci <branch> <pull_request_number|false>"
  echo
  echo "Default build/push tag is current git branch ($git_branch)"
}

# Parse command line arguments and do work. Should be called like so:
#   main $*
main ()
{
  # Figure out which tag to use for this image build; either use the name passed
  # as $2, or derive a tag from the branch name using a simple heuristic.
  if [ -n "$2" ]
  then
    git_branch=$2
    echo "Assuming Git branch '$git_branch' by explicit request"
  else
    git_ref=$(git symbolic-ref HEAD 2>/dev/null)
    git_branch=${git_ref##refs/heads/}
    echo "Detected Git branch '$git_branch'"
  fi

  # Map the git branch name to an image tag
  case $git_branch in
    master)
      tag=latest
      ;;
    production)
      tag="${git_branch}-isolated"
      ;;
    *)
      tag=$git_branch
      ;;
  esac
  echo "Derived Docker image tag '$tag' from Git branch name"

  pull_request_number=$3

  # Run the command
  case $1 in
  ci)
    ci $tag $pull_request_number
    ;;
  build)
    build $tag
    ;;
  clean)
    clean
    ;;
  push)
    push $tag
    ;;
  *)
    help
    ;;
  esac
}

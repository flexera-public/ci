# VARIABLES
DEFAULT_APT_DOCKER_PKG="docker-ce=18.06.3~ce~3-0~ubuntu"

docker_install()
{
  set -e # Automatically exit on error

  if  [ -z "$APT_DOCKER_PKG" ]
  then
    echo "*** APT_DOCKER_PKG undefined, using default: $DEFAULT_APT_DOCKER_PKG"
    export APT_DOCKER_PKG=$DEFAULT_APT_DOCKER_PKG
  fi

  echo "*** Adding get.docker.com repository"
  sudo curl -sSL https://get.docker.com/ | sh
  echo "*** Installing $APT_DOCKER_PKG"
  sudo apt-get -y --allow-downgrades $APT_PARAMS -o Dpkg::Options::="--force-confnew" install $APT_DOCKER_PKG && docker -v
}

check_sudo()
{
  if [[ $TRAVIS_SUDO != "true" ]]
  then
    echo "!!! ERROR: you need a sudo-enabled (sudo: required) Travis build to use the specified options in travis-setup.sh script (probably you are requesting Docker)"
    exit 10
  fi
}


function build_status {
  if [ -z "${TRAVIS_PRO_TOKEN}" ]
  then
    echo "ERROR!!!!! Need TRAVIS_PRO_TOKEN to get build status"
    exit 10
  fi

  repo=$1
  sha=$2
  curl -s -H "Authorization: token $TRAVIS_PRO_TOKEN" https://api.travis-ci.com/repos/rightscale/${repo}/builds | jq ".[] | {sha: .commit, state: .state, result: .result, branch: .branch} | select(.sha == \"${sha}\" and .state == \"finished\" )"
}

function build_already_green {
  if [[ $CI_DISABLE_SHA_GREEN_CHECK =~ ^(true|TRUE|1)$ ]]
  then
    echo "WARNING!!! CI_DISABLE_GREEN_CHECK is set, so always running tests"
    return 1
  fi

  if [ -z "${TRAVIS_PRO_TOKEN}" ]
  then
    echo "WARNING!!!!! Need TRAVIS_PRO_TOKEN defined to be able to use green-SHA-skip-tests feature"
    return 1 # false
  fi

  repo=$1
  sha=$2
  res=$(build_status $repo $sha)
  if [[ "0${res}" != "0" ]] ; then
    state=$(echo $res | jq '.state' | tr -d '"')
    result=$(echo $res | jq '.result' | tr -d '"')
    if [[ "0$state" == "0finished" && "0$result" == "00" ]] ; then
      echo "***********************************************************************************"
      echo "***********************************************************************************"
      echo " WARNING!! Skipping tests since SHA $sha has already been tested with green status"
      echo "***********************************************************************************"
      echo "***********************************************************************************"
      return 0 # true
    fi
  fi
  return 1 # false
}


if [[ $DOCKER =~ ^(true|TRUE|1)$ ]]
then
  check_sudo
  docker_install
fi

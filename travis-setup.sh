# VARIABLES
DEFAULT_APT_DOCKER_PKG="docker-engine=1.12.5-0~ubuntu-trusty"

docker_install()
{
  set -e # Automatically exit on error

  if  [ -z "$APT_DOCKER_PKG" ]
  then
    echo "*** APT_DOCKER_PKG undefined, using default: $DEFAULT_APT_DOCKER_PKG"
    export APT_DOCKER_PKG=$DEFAULT_APT_DOCKER_PKG
  fi

  if [[ $APT_ALLOW_UNAUTHENTICATED == "true" ]]
  then
    echo "WARNING!!!! Installing unauthenticated packages !!!!"
    export APT_PARAMS="--allow-unauthenticated"
  else
    export APT_PARAMS=""
  fi
  echo "*** Adding download.docker.com repository (trusty)"
  echo 'deb "https://download.docker.com/repo" ubuntu-trusty main' | sudo tee /etc/apt/sources.list.d/docker.list
  sudo apt-get update
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

install_docker_repo_gpg_key()
{
  # this is the public key for f76221572c52609d, which can be found at
  # http://keyserver.ubuntu.com/pks/lookup?op=vindex&search=0xF76221572C52609D&fingerprint=on
  sudo tee ./2C52609D.pub.key > /dev/null << 'EOF'
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: SKS 1.1.6
Comment: Hostname: keyserver.ubuntu.com

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
=lwu7
-----END PGP PUBLIC KEY BLOCK-----

EOF

  # install the key and delete the file.
  sudo apt-key add ./2C52609D.pub.key
  sudo rm -f ./2C52609D.pub.key
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
  install_docker_repo_gpg_key
  docker_install
fi

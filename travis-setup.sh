# Automatically exit on error
set -e

# VARIABLES
DEFAULT_APT_DOCKER_PKG="docker-engine=1.12.5-0~ubuntu-trusty"

docker_install()
{
  if  [ -z "$APT_DOCKER_PKG" ]
  then
    echo "*** APT_DOCKER_PKG undefined, using default: $DEFAULT_APT_DOCKER_PKG"
    export APT_DOCKER_PKG=$DEFAULT_APT_DOCKER_PKG
  fi

  echo "*** Adding dockerproject key"
  sudo apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys F76221572C52609D
  echo "*** Adding apt.dockerproject.org repository (trusty)"
  echo 'deb "https://apt.dockerproject.org/repo" ubuntu-trusty main' | sudo tee /etc/apt/sources.list.d/docker.list
  sudo apt-get update
  echo "*** Installing $APT_DOCKER_PKG"
  sudo apt-get -y --allow-downgrades -o Dpkg::Options::="--force-confnew" install $APT_DOCKER_PKG && docker -v
}

check_sudo()
{
  if [[ $TRAVIS_SUDO != "true" ]]
  then
    echo "!!! ERROR: you need a sudo-enabled (sudo: required) Travis build to use the specified options in travis-setup.sh script (probably you are requesting Docker)"
    exit 10
  fi
}


if [[ $DOCKER =~ ^(true|TRUE|1)$ ]]
then
  check_sudo
  docker_install
fi

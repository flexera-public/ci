function check_sudo {
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

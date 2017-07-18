# ci
A collection of public resources used by RightScale CI

## travis-setup.sh

A script designed to help setting up the build environment.

Download and call it at the beginning of the "before_install" section of your build:

```
before_install:
  - curl -s https://raw.githubusercontent.com/rightscale/ci/TL-1013_travis-setup/travis-setup.sh | DOCKER=true bash
```

Pass options/commands via environment variables, these are the available options:

- `DOCKER`: if set to **true**, it will install the specified docker version in APT_DOCKER_PKG (if APT_DOCKER_PKG is unset, it will install the script's default)


## docker-shared.sh

A script designed to be retrieved and executed from CI builds to manage the building and pushing of your containers.

Needs the following environment variables:

- `DOCKERHUB_USER`: DockerHub username
- `DOCKERHUB_PASSWORD`: DockerHub password

Provides the following functions:

- `clean`: Remove intermediate build artifacts
- `login`: logs into DockerHub
- `build`: builds Docker image
- `push`: Pushes Docker image
- `ci(branch, is_pull_request)`: Depending on the branch name, the pull request status and other conditions that get retrieved at execution time, this function ends executing the `login`, `build` and `push` functions

## Maintained by
 - [Christian Teijon](https://github.com/crunis)

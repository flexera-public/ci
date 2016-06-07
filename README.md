# ci
A collection of public resources used by RightScale CI

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

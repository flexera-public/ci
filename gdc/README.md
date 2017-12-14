# GDC - Go Dependency Checker 

*gdc* is able to read your Go project imports and process Git changes between two hashes and tell you if your service needs to be tested and rebuild. 

This is useful for Git monorepos with several Go services to avoid a new commit to trigger a new cycle of test + build the services that have not changed.

## How to use in Travis

If the services in your repo live in separate directories, that are specified as part of a Travis "matrix build", like for example:

```yaml
env:
  matrix:
    - SERVICE=queue
    - SERVICE=billing
    - SERVICE=image/shrinker
    - SERVICE=image/recogniser
```

You can add these lines in your .travis.yml file so only the services with modified files will be processed in TravisCI:

```yaml
before_install:
  - curl https://raw.githubusercontent.com/rightscale/ci/v1/gdc/bin/gdc_linux -o gdc_linux && chmod a+x ./gdc_linux
  - export DEPS=$(./gdc_linux travis $SERVICE)
  - if [[ "$DEPS" == "skip" ]]; then echo "Skipping $SERVICE since no dependencies changed"; travis_terminate 0; else echo "Hit dependencies $DEPS"; fi
```

## GDC commands

gdc provides several commands that can be handy in the command line

### version

```bash
gdc version
```

Returns gdc version

### gitdiff

```bash
gdc -sha1 <sha1> -sha2 <sha2> gitdiff
```

Lists all the changed files between 2 commits. Note:

### imports

```bash
gdc imports <directory>
```
Will list all the "import clauses" of a given directory

### deps

```bash
gdc deps <directory>
```
Will list all the dependencies of a given directory, that is, all files or directories that should trigger a rebuild of the corresponding service if any of these files changed. This is the reason the files of the directory are also included.

### check

```bash
gdc -sha1 <sha1> -sha2 <sha2> check <directory>
```

Shows, given a commit range and a directory, the dependencies that got modified.

### travis

```bash
gdc travis <directory>
```

This command is almost identical to the check command explained in the previous section but with some variations that make it handier to run in a Travis build. They are:

- sha1 and sha2 are extracted from the environment variable TRAVIS_COMMIT_RANGE, if this var is empty, it defaults to HEAD + HEAD~1 (see https://docs.travis-ci.com/user/environment-variables/#Default-Environment-Variables)
- It will return the string "skip" if there are no dependencies hit, otherwise it will return a message with the dependencies.

## Notes

- If any file in the root directory of the project changes, that will be considered a dependency. Unfortunately this includes changes to README.md, etc.. But covers for changes on glide.yaml, ...
- Any file change inside the given directory will be considered a dependency
- It doesn't matter if sha1 is older or newer than sha2, the output is always the same, i.e., swapping sha1 and sha2 produces the same result
- sha1 and sha2 can be specified with full SHAs (40 characters), shorter SHAs of any size (as long as they are unambiguous) or HEAD~X references.
- Like "git" command, gdc will try to find a git project in the current directory and travel up the directory hierarchy until it finds it.

## ToDo

- Currently, GDC is not able to recursively account for depedencies, that is, if service A depends on service B code, and service B depends on service C code, a change on service C code, won't be detected as a dependency of service A (but it will be detected as a dependency of service B).

## HowTo build for diferent architectures

MacOs

```bash
GOOS=darwin GOARCH=amd64 go build -o bin/gdc_mac *.go
```

Linux

```bash
GOOS=linux GOARCH=amd64 go build -o bin/gdc_linux *.go
```

Windows

```bash
GOOS=windows GOARCH=386 go build -o bin/gdc_windows.exe *.go
```

## License

The `gdc` source code is subject to the MIT license, see the
[LICENSE](https://github.com/rightscale/ci/gdc/LICENSE) file.

// The MIT License (MIT)
//
// Copyright (c) 2017 RightScale, Inc, All Rights Reserved Worldwide.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Verbose : If true, display debug messages
var Verbose = false
var version = "0.1.0"

// Returns flags and params from command line
func getFlagsAndParams() (flags map[string]string, command string, directory string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:  %s [flags] <command> [directory] \n", os.Args[0], os.Args[0])
		fmt.Println("\nAvailable commands:\n")
		fmt.Println("  version - returns current version")
		fmt.Println("  check - check if directory has changed dependencies")
		fmt.Println("  travis - check if directory has changed dependencies, using Travis Env Vars")
		fmt.Println("  gitdiff - show changed files")
		fmt.Println("  root - show root directories that have changed dependencies")
		fmt.Println("  deps - show all deps of a given directory (this will include the files of the directory)")
		fmt.Println("  imports - show all imports of a given directory")
		fmt.Println("\nAvailable flags:\n")
		flag.PrintDefaults()
		fmt.Println(" ")
	}

	verbose := flag.Bool("verbose", false, "enable verbose mode")
	usetravisenv := flag.Bool("usetravisenv", false, "use TRAVIS_COMMIT_RANGE env var")
	sha1 := flag.String("sha1", "HEAD", "sha1, defaults to HEAD")
	sha2 := flag.String("sha2", "HEAD~1", "sha2, defaults to HEAD~1")
	flag.Parse()

	flags = make(map[string]string)
	flags["sha1"] = *sha1
	flags["sha2"] = *sha2
	flags["usetravisenv"] = strconv.FormatBool(*usetravisenv)
	Verbose = *verbose

	params := os.Args[len(os.Args)-flag.NArg() : len(os.Args)]

	if len(params) == 0 {
		fmt.Println("You must specify a command")
		flag.Usage()
		os.Exit(1)
	}
	command = params[0]
	directory = ""
	if len(params) == 2 {
		directory = params[1]
	}

	return
}

// Returns env var GOPATH
func getGoPath() (gopath string) {
	return os.Getenv("GOPATH")
}

// Returns current project path, relative to GOPATH/src
// Obtained in this way: (GIT PATH) - (GOPATH)
//    Also removing /src preffix and .git/ suffix
func getCurrentRelativePath() (relPath string) {
	workdir := getRepoPath() // Finds current Git repo base path
	workdir = strings.TrimSuffix(workdir, ".git/")
	gopath := getGoPath()
	if !strings.HasPrefix(workdir, gopath) {
		fmt.Printf("ERROR! Current working repository dir should be inside GOPATH.\n")
		fmt.Printf("       Current repository directory: %s\n", workdir)
		fmt.Printf("       GOPATH: %s\n", gopath)
		os.Exit(1)
	}
	relPath = strings.TrimPrefix(workdir, gopath)
	relPath = strings.TrimPrefix(relPath, "/src/")

	return
}

// Returns sha1 and sha2 from env var TRAVIS_COMMIT_RANGE
func getTravisCommitRange() (sha1 string, sha2 string) {
	commitRange := os.Getenv("TRAVIS_COMMIT_RANGE")
	if Verbose && commitRange != "" {
		fmt.Printf("TRAVIS_COMMIT_RANGE %v\n", commitRange)
	}

	if len(commitRange) == 0 {
		if Verbose {
			fmt.Println("Empty TRAVIS_COMMIT_RANGE, assuming HEAD~1..HEAD")
		}
		sha1 = "HEAD~1"
		sha2 = "HEAD"
		return
	}

	shas := strings.Split(commitRange, "...")
	if len(shas) != 2 {
		errTxt := fmt.Sprintf("Error parsing commit Range(%v). %v SHAs found", commitRange, len(shas))
		panic(errors.New(errTxt))
	}

	return shas[0], shas[1]
}

// Given a list of imports and the Git changed paths, returns an array of hit dependencies
//   If any of the changedPaths is a root file, it always counts as dependency
func hitDepends(imports, changedPaths []string) (depends []string) {
	for _, path := range changedPaths {
		if isRootFile(path) {
			depends = append(depends, path)
			continue
		}
		for _, anImport := range imports {
			if anImport == filepath.Dir(path) {
				depends = append(depends, anImport)
				break
			}
		}
	}

	if Verbose == true {
		fmt.Println("IMPORTS =========================")
		for _, anImport := range imports {
			fmt.Printf("  %s \n", anImport)
		}
		fmt.Println("PATHS ===========================")
		for _, path := range changedPaths {
			fmt.Printf("  %s \n", path)
		}
		fmt.Println("HITS ===========================")
		for _, hit := range depends {
			fmt.Printf("  %s \n", hit)
		}
	}

	return
}

// Given SHA1 and SHA2 and a directory, checks if there are hit dependencies
func findHitDeps(sha1, sha2, directory string) []string {
	paths := changedPaths(sha1, sha2)
	imports := getParsedDependencies(directory, getCurrentRelativePath())

	return hitDepends(imports, paths)
}

// Outputs the Git modified files between sha1 and sha2
func showGitDiff(sha1, sha2 string) {
	if sha1 == "" || sha2 == "" {
		fmt.Println("Using TRAVIS_COMMIT_RANGE for sha1 and sha2")
		sha1, sha2 = getTravisCommitRange()
	}
	paths := changedPaths(sha1, sha2)
	fmt.Println("Changed Paths folders:")
	for _, path := range paths {
		fmt.Println(path)
	}

	return
}

// Travis functionality, returns dependencies or "skip" if there are no hit dependencies
func travis(directory string) {
	sha1, sha2 := getTravisCommitRange()

	depends := findHitDeps(sha1, sha2, directory)
	if len(depends) == 0 {
		fmt.Printf("skip\n")
		os.Exit(0)
	} else {
		fmt.Printf("%v\n", depends)
		os.Exit(0)
	}
}

func main() {
	var sha1, sha2 string
	flags, command, directory := getFlagsAndParams()
	sha1 = flags["sha1"]
	sha2 = flags["sha2"]
	if flags["usetravisenv"] == "true" {
		fmt.Println("Use travis env activated")
		sha1, sha2 = getTravisCommitRange()
	}

	if Verbose {
		fmt.Printf("Current repo path: %s \n", getRepoPath())
		fmt.Printf("Current GO path: %s \n", getGoPath())
		fmt.Printf("Current relative path: %s \n", getCurrentRelativePath())
		fmt.Printf("SHA1: %s  SHA2: %s \n", sha1, sha2)
	}

	switch command {
	case "version":
		fmt.Printf("gdc version %s\n", version)
	case "travis":
		if len(directory) == 0 {
			fmt.Println("You need to specify a directory with this command")
			os.Exit(1)
		}
		travis(directory)
	case "root":
		sha1, sha2 := getTravisCommitRange()
		folders := changedRootFolders(sha1, sha2)
		fmt.Printf("Changed ROOT folders: %v\n", folders)
	case "deps":
		if len(directory) == 0 {
			fmt.Println("You need to specify a directory with this command")
			os.Exit(1)
		}
		res := getParsedDependencies(directory, getCurrentRelativePath())
		fmt.Printf("Parsed dependencies on directory %s: \n\n", directory)
		for _, k := range res {
			fmt.Println(k)
		}
	case "imports":
		if len(directory) == 0 {
			fmt.Println("You need to specify a directory with this command")
			os.Exit(1)
		}
		res := getImports(directory)
		fmt.Printf("Dependencies on directory %s: \n\n", directory)
		for _, k := range res {
			fmt.Println(k)
		}
	case "gitdiff":
		showGitDiff(sha1, sha2)
	case "check":
		if len(directory) == 0 {
			fmt.Println("You need to specify a directory with this command")
			os.Exit(1)
		}
		depends := findHitDeps(sha1, sha2, directory)
		if len(depends) > 0 {
			fmt.Printf("Dependencies found: %v \n", depends)
		} else {
			fmt.Println("No dependencies found")
		}
	default:
		fmt.Printf("Unknown command %s\n", command)
	}
}

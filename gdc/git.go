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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var gitPath = ""

func getRepoPath() string {
	if len(gitPath) == 0 {
		openRepo()
	}
	return gitPath
}

func changedPaths(sh1, sh2 string) []string {
	if Verbose {
		fmt.Printf("changedPaths from %s to %s \n", sh1, sh2)
	}
	sha1 := plumbing.NewHash(expandSHA(sh1))
	sha2 := plumbing.NewHash(expandSHA(sh2))
	repo := openRepo()
	commit1, err := repo.CommitObject(sha1)
	if err != nil {
		panic(err)
	}

	commit2, err := repo.CommitObject(sha2)
	if err != nil {
		panic(err)
	}

	patch, err := commit1.Patch(commit2)
	if err != nil {
		panic(err)
	}

	paths := make(map[string]struct{})
	fpatches := patch.FilePatches()
	for _, fpatch := range fpatches {
		from, to := fpatch.Files()
		if from != nil {
			paths[from.Path()] = struct{}{}
		}
		if to != nil {
			paths[to.Path()] = struct{}{}
		}
	}

	return getSortedKeys(paths)
}

func changedDirs(sh1, sh2 string) []string {
	if Verbose {
		fmt.Printf("changedDirs from %s to %s \n", sh1, sh2)
	}
	sha1 := plumbing.NewHash(expandSHA(sh1))
	sha2 := plumbing.NewHash(expandSHA(sh2))
	repo := openRepo()
	commit1, err := repo.CommitObject(sha1)
	if err != nil {
		panic(err)
	}

	commit2, err := repo.CommitObject(sha2)
	if err != nil {
		panic(err)
	}

	patch, err := commit1.Patch(commit2)
	if err != nil {
		panic(err)
	}

	paths := make(map[string]struct{})
	fpatches := patch.FilePatches()
	for _, fpatch := range fpatches {
		from, to := fpatch.Files()
		if from != nil {
			paths[filepath.Dir(from.Path())] = struct{}{}
		}
		if to != nil {
			paths[filepath.Dir(to.Path())] = struct{}{}
		}
	}

	return getSortedKeys(paths)
}

func listBranches(repo *git.Repository) {
	branches, _ := repo.Branches()
	branch, _ := branches.Next()

	var err error
	for branch != nil && err == nil {
		fmt.Printf("%v \n", branch.Name())
		branch, err = branches.Next()
	}
}

func listReferences() {
	repo := openRepo()
	refs, _ := repo.References()
	ref, _ := refs.Next()

	var err error
	for ref != nil && err == nil {
		fmt.Printf("%v - %v \n", ref.Name(), ref.Hash())
		ref, err = refs.Next()
	}
}

func listObjects() {
	repo := openRepo()
	objs, _ := repo.Objects()
	obj, _ := objs.Next()

	var err error
	for obj != nil && err == nil {
		fmt.Printf("%v - %v \n", obj.Type(), obj.ID())
		obj, err = objs.Next()
	}
}

func commitInfo(commit *object.Commit) {
	fmt.Printf("%v [%v] - %v \n", commit.Hash, commit.NumParents(), commit.Message)
	fmt.Printf("   TREEHASH: %v \n", commit.TreeHash)

	for _, phash := range commit.ParentHashes {
		fmt.Printf("     PH: %v \n", phash)
	}

	fmt.Println()
}

func listCommits() {
	var err error
	repo := openRepo()
	commits, err := repo.CommitObjects()
	if err != nil {
		panic(err)
	}
	commit, err := commits.Next()
	if err != nil {
		panic(err)
	}

	for commit != nil && err == nil {
		commitInfo(commit)
		commit, err = commits.Next()
	}
}

func expandSHA(givenSHA string) (realSHA string) {
	switch {
	case strings.HasPrefix(givenSHA, "HEAD"):
		return translateHead(givenSHA)
	case len(givenSHA) == 40:
		return givenSHA
	case len(givenSHA) > 40:
		errTxt := fmt.Sprintf("Error, the given SHA (%v), is longer than 40 (length is %v)\n", givenSHA, len(givenSHA))
		panic(errors.New(errTxt))
	default:
		return getFullSHA(givenSHA)
	}
}

func getFullSHA(shortSHA string) (fullSHA string) {
	repo := openRepo()
	commitIter, err := repo.CommitObjects()
	if err != nil {
		panic(err)
	}

	var results []string
	err = commitIter.ForEach(func(commit *object.Commit) error {
		if strings.HasPrefix(commit.Hash.String(), shortSHA) {
			results = append(results, commit.Hash.String())
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	switch n := len(results); n {
	case 0:
		errTxt := fmt.Sprintf("No SHA found starting with %s", shortSHA)
		panic(errors.New(errTxt))
	case 1:
		fullSHA = results[0]
	default:
		errTxt := fmt.Sprintf("short SHA %s is ambiguous, %d matching SHAs found", shortSHA, n)
		panic(errors.New(errTxt))
	}

	return
}

func openRepo() *git.Repository {
	startDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Keep looking for a git repo in the parent dir
	repo, err := git.PlainOpen(startDir)
	for err != nil && err.Error() == "repository not exists" {
		var prevDir string
		startDir, prevDir = filepath.Join(startDir, ".."), startDir
		if startDir == prevDir {
			panic(errors.New("Can't find git repo in current dir or any of its  parents\n"))
		}
		repo, err = git.PlainOpen(startDir)
	}
	gitPath = startDir

	return repo
}

func stripFilenames(fullPaths []string) []string {
	dirNames := make(map[string]struct{})
	for _, fullPath := range fullPaths {
		dirName := filepath.Dir(fullPath)
		dirNames[dirName] = struct{}{}
	}
	return getSortedKeys(dirNames)
}

func headExtractN(headTag string) (n int) {
	if strings.Compare(headTag, "HEAD") == 0 {
		return 0
	}
	if !strings.HasPrefix(headTag, "HEAD~") {
		errStr := fmt.Sprintf("%s git ref should start with HEAD~", headTag)
		panic(errors.New(errStr))
	}
	res := strings.Trim(headTag, "HEAD~")

	n, err := strconv.Atoi(res)
	if err != nil {
		panic(err)
	}
	return n
}

func translateHead(headTag string) (sha string) {
	var err error
	repo := openRepo()
	headRef, err := repo.Head()
	if err != nil {
		panic(err)
	}
	commit, err := repo.CommitObject(headRef.Hash())

	for i := 0; i < headExtractN(headTag); i++ {
		commit, err = repo.CommitObject(commit.ParentHashes[0])
		if Verbose {
			fmt.Printf(" %s  \n", commit.Hash.String())
		}
		if err != nil {
			panic(err)
		}
	}

	sha = commit.Hash.String()

	if Verbose {
		fmt.Printf(" %s translated to %s \n", headTag, sha)
	}

	return sha
}

func extractDirNames(fullPaths []string) []string {
	dirNames := make(map[string]struct{})
	for _, fullPath := range fullPaths {
		dirName := filepath.Dir(fullPath)
		dirNames[dirName] = struct{}{}
	}
	return getSortedKeys(dirNames)
}

func changedRootFolders(sha1, sha2 string) (rootFolders []string) {
	paths := changedPaths(sha1, sha2)

	for _, path := range paths {
		if !strings.ContainsAny(path, "/") {
			path = "ROOT"
		}
		folders := strings.Split(path, "/")
		if !contains(rootFolders, folders[0]) {
			rootFolders = append(rootFolders, folders[0])
		}
	}

	return
}

func contains(s []string, e string) bool {
	for _, x := range s {
		if x == e {
			return true
		}
	}
	return false
}

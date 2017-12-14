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
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func printImports(path string) {
	// Print the imports from the file's AST.
	for _, s := range getFileImports(path) {
		fmt.Println(s)
	}
}

func getFileImports(path string) []string {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		log.Fatal(err)
	}

	imports := make([]string, len(f.Imports))
	for i, s := range f.Imports {
		imports[i] = s.Path.Value
	}
	return imports
}

func getGoFiles(path string) []string {
	fileList := []string{}
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return fileList
}

func getAllFiles(path string) []string {
	fileList := []string{}
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return fileList
}

func getExecutionDir() string {
	// currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return currentDir
}

func getSortedKeys(aMap map[string]struct{}) []string {
	keys := make([]string, len(aMap))

	i := 0
	for k := range aMap {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func getImports(directory string) (files []string) {
	imports := make(map[string]struct{}) // struct{} occupies 0 bytes

	for _, file := range getGoFiles(directory) {
		for _, impoort := range getFileImports(file) {
			imports[impoort] = struct{}{}
		}
	}

	files = getSortedKeys(imports)
	return
}

func getParsedImports(directory string, projectDir string) (imports []string) {
	for _, impoort := range getImports(directory) {
		impoort = strings.Trim(impoort, "\"")
		if strings.HasPrefix(impoort, projectDir) {
			impoort = strings.TrimPrefix(impoort, projectDir)
			impoort = strings.TrimPrefix(impoort, "/")
			imports = append(imports, impoort)
		}
	}

	return
}

func getParsedDependencies(directory string, projectDir string) (imports []string) {
	// Adds filenames to imports (for non-Go files)
	deps := make(map[string]struct{})

	for _, anImport := range getParsedImports(directory, projectDir) {
		deps[anImport] = struct{}{}
	}
	for _, file := range getAllFiles(directory) {
		deps[file] = struct{}{}
	}

	return getSortedKeys(deps)
}

func isRootFile(path string) bool {
	switch filepath.Dir(path) {
	case ".", "/":
		return true
	default:
		return false
	}
}

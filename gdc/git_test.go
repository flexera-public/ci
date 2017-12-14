package main

import (
	"testing"
	"reflect"
)

func TestHeadExtractN(t *testing.T){
	var testString = "HEAD"
	n := headExtractN(testString)
	if n != 0 {
		t.Error(testString, "should return 0, got", n)
	}

	testString = "HEAD~1"
	n = headExtractN(testString)
	if n != 1 {
		t.Error(testString, "should return 1, got", n)
	}

	testString = "HEAD~5"
	n = headExtractN(testString)
	if n != 5 {
		t.Error(testString, "should return 5, got", n)
	}

	testString = "HEAD~13"
	n = headExtractN(testString)
	if n != 13 {
		t.Error(testString, "should return 13, got", n)
	}

	testString = "HEAD~101"
	n = headExtractN(testString)
	if n != 101 {
		t.Error(testString, "should return 101, got", n)
	}

	testPanic(t)
}

func testPanic(t *testing.T){
	testString := "HEAD3"

	defer func() {
		if r := recover(); r == nil {
				t.Errorf("The code did not panic with %s", testString)
		}
	}()

	_ = headExtractN(testString)
}

func TestExtractDirnames(t *testing.T){
	test := []string{ "uno/dos/tres.go", "uno/dos/tres/cuatro.go", "uno/dos/cuatro.go" }
	expected := []string { "uno/dos", "uno/dos/tres" }
	res := extractDirNames(test)
	if ! reflect.DeepEqual(res, expected) {
		t.Error(test, "should return", expected, "but returned", res)
	}

}

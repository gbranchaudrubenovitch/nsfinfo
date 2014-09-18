package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// tests here are aimed against the methods in nsfinfo.go

func TestGetHeaderWithFileNotFound(t *testing.T) {
	testGetHeader(t, "path-to-no-file", "open path-to-no-file", func(c string) (*nsfHeader, error) { return getHeader(c) })
}

func TestGetHeaderWithFileTooShort(t *testing.T) {
	testGetHeader(t, "less than 128 bytes", "unexpected EOF", getHeaderFromSlice)
}

func TestGetHeaderWithInvalidHeader(t *testing.T) {
	testGetHeader(t, "12812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812812", "invalid FourCC", getHeaderFromSlice)
}

// testing main() is really about adding code coverage
func TestMainWithValidFile(t *testing.T) {
	oldStdout := os.Stdout
	os.Stdout = nil
	defer func() { os.Stdout = oldStdout }() // not really interested in seeing main()'s stdout, so hide it for the test

	os.Args = []string{"nsfinfo", "samples/smb1.nsf"}
	main()
}

var getHeaderFromSlice = func(c string) (*nsfHeader, error) {
	ioutil.WriteFile("t.nsf", []byte(c), 0644)
	defer os.Remove("t.nsf")
	return getHeader("t.nsf")
}

func testGetHeader(t *testing.T, content string, expectedError string, getHeader func(string) (*nsfHeader, error)) {
	h, e := getHeader(content)
	hasErrorNoHeader(t, h, e)

	mustContain(t, e.Error(), expectedError)
}

func hasErrorNoHeader(t *testing.T, h *nsfHeader, e error) {
	if h != nil || e == nil {
		t.Fatalf("header must be nil & error must be non-nil (h is nil: %v | e is nil: %v).", h == nil, e == nil)
	}
}

func hasHeaderNoError(t *testing.T, h *nsfHeader, e error) {
	if h == nil || e != nil {
		t.Fatalf("header must not be nil & error must be nil (h is nil: %v | e is nil: %v).", h == nil, e == nil)
	}
}

func mustContain(t *testing.T, hay string, needle string) {
	if !strings.Contains(hay, needle) {
		t.Fatalf("string \"%v\"\n\tshould contain \"%s\"", hay, needle)
	}
}

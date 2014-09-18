package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// tests here are aimed against the methods in nsfinfo.go

func TestGetHeaderWithFileNotFound(t *testing.T) {
	h, e := getHeader("path-to-no-file")
	hasErrorNoHeader(t, h, e)

	mustContain(t, e.Error(), "open path-to-no-file")
}

func TestGetHeaderWithFileTooShort(t *testing.T) {
	h, e := getHeaderFromSlice([]byte("less than 128 bytes"))
	hasErrorNoHeader(t, h, e)

	mustContain(t, e.Error(), "unexpected EOF")
}

func TestGetHeaderWithInvalidHeader(t *testing.T) {
	h, e := getHeaderFromSlice([]byte("11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"))
	hasErrorNoHeader(t, h, e)

	mustContain(t, e.Error(), "invalid FourCC")
}

// testing main() is really about adding code coverage
func TestMainWithValidFile(t *testing.T) {
	oldStdout := os.Stdout
	os.Stdout = nil
	defer func() { os.Stdout = oldStdout }() // not really interested in seeing main()'s stdout, so hide it for the test

	os.Args = []string{"nsfinfo", "samples/smb1.nsf"}
	main()
}

// getHeaderFromSlice calls getHeader() with the path to a temp file created from the byte slice
func getHeaderFromSlice(content []byte) (*nsfHeader, error) {
	ioutil.WriteFile("t.nsf", content, 0644)
	defer os.Remove("t.nsf")
	return getHeader("t.nsf")
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

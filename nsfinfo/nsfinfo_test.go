package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestGetHeaderWithFileNotFound(t *testing.T) {
	h, e := getHeader("path-to-no-file")
	nilHeaderNonNilError(t, h, e)

	if !strings.Contains(e.Error(), "open") {
		t.Errorf("have %+v\n\twant %+v", e, "error on open() call")
	}
}

func TestGetHeaderWithFileTooShort(t *testing.T) {
	fileName := "fileTooShort.nsf"
	e := ioutil.WriteFile(fileName, []byte("less than 128 bytes"), 0644)
	defer os.Remove(fileName)
	if e != nil {
		panic(e)
	}

	h, e := getHeader(fileName)
	nilHeaderNonNilError(t, h, e)

	if !strings.Contains(e.Error(), "unexpected EOF") {
		t.Errorf("have %+v\n\twant %+v", e, "unexpected EOF (file is too short)")
	}
}

func nilHeaderNonNilError(t *testing.T, h *nsfHeader, e error) {
	if h != nil || e == nil {
		t.Errorf("header must be nill & error must be non-nill.")
	}
}

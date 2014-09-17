package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestGetHeaderWithFileNotFound(t *testing.T) {
	h, e := getHeader("path-to-no-file")
	hasErrorNoHeader(t, h, e)

	if !strings.Contains(e.Error(), "open") {
		t.Fatalf("have %+v\n\twant %+v", e, "error on open() call")
	}
}

func TestGetHeaderWithFileTooShort(t *testing.T) {
	fileName := "file-too-short.nsf"
	ioutil.WriteFile(fileName, []byte("less than 128 bytes"), 0644)
	defer os.Remove(fileName)

	h, e := getHeader(fileName)
	hasErrorNoHeader(t, h, e)

	if !strings.Contains(e.Error(), "unexpected EOF") {
		t.Fatalf("have %+v\n\twant %+v", e, "unexpected EOF (file is too short)")
	}
}

var validNsfHeader = []byte{
	0x4E, 0x45, 0x53, 0x4D, 0x1A, 0x01, 0x08, 0x01, 0x60, 0x8D, 0x03, 0xA0, 0x00, 0xA0, 0x54, 0x68,
	0x65, 0x20, 0x4C, 0x65, 0x67, 0x65, 0x6E, 0x64, 0x20, 0x6F, 0x66, 0x20, 0x5A, 0x65, 0x6C, 0x64,
	0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4B, 0x6F,
	0x6A, 0x69, 0x20, 0x4B, 0x6F, 0x6E, 0x64, 0x6F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x31, 0x39,
	0x38, 0x37, 0x20, 0x4E, 0x69, 0x6E, 0x74, 0x65, 0x6E, 0x64, 0x6F, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1A, 0x41,
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

func TestGetHeaderWithInvalidNsfPrelude(t *testing.T) {
	var invalidPrelude = make([]byte, len(validNsfHeader))
	copy(invalidPrelude, validNsfHeader)
	invalidPrelude[1] = 0x55

	fileName := "invalid-prelude.nsf"
	ioutil.WriteFile(fileName, invalidPrelude, 0644)
	defer os.Remove(fileName)

	h, e := getHeader(fileName)
	hasErrorNoHeader(t, h, e)

	if !strings.Contains(e.Error(), "invalid prelude") {
		t.Fatalf("have %+v\n\twant %+v", e, "error about invalid prelude")
	}
}

func TestGetHeaderWithInvalidExpansionChips(t *testing.T) {
	var invalidExpansion = make([]byte, len(validNsfHeader))
	copy(invalidExpansion, validNsfHeader)
	invalidExpansion[123] = 1 << 7

	fileName := "invalid-expansion.nsf"
	ioutil.WriteFile(fileName, invalidExpansion, 0644)
	defer os.Remove(fileName)

	h, e := getHeader(fileName)
	hasErrorNoHeader(t, h, e)

	if !strings.Contains(e.Error(), "extra sound chip section") {
		t.Fatalf("have %+v\n\twant %+v", e, "error about expansion section")
	}
}

func hasErrorNoHeader(t *testing.T, h *nsfHeader, e error) {
	if h != nil || e == nil {
		t.Fatalf("header must be nill & error must be non-nil (h is nil: %v | e is nil: %v).", h == nil, e == nil)
	}
}

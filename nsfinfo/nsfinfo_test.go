package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var validHeader = []byte{
	0x4E, 0x45, 0x53, 0x4D, 0x1A, 0x01, 0x08, 0x01, 0x60, 0x8D, 0x03, 0xA0, 0x00, 0xA0, 0x54, 0x68,
	0x65, 0x20, 0x4C, 0x65, 0x67, 0x65, 0x6E, 0x64, 0x20, 0x6F, 0x66, 0x20, 0x5A, 0x65, 0x6C, 0x64,
	0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4B, 0x6F,
	0x6A, 0x69, 0x20, 0x4B, 0x6F, 0x6E, 0x64, 0x6F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x31, 0x39,
	0x38, 0x37, 0x20, 0x4E, 0x69, 0x6E, 0x74, 0x65, 0x6E, 0x64, 0x6F, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1A, 0x41,
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

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

func TestGetHeaderWithInvalidNsfPrelude(t *testing.T) {
	var invalidPrelude = append([]byte(nil), validHeader...)
	invalidPrelude[1] = 0x55 // scraps the 2nd prelude byte

	h, e := getHeaderFromSlice(invalidPrelude)
	hasErrorNoHeader(t, h, e)

	mustContain(t, e.Error(), "invalid prelude")
}

func TestGetHeaderWithInvalidExpansionChips(t *testing.T) {
	var invalidExpansion = append([]byte(nil), validHeader...)
	invalidExpansion[123] = 1 << 7 // this is one of the "future chip" bits that must never be set

	h, e := getHeaderFromSlice(invalidExpansion)
	hasErrorNoHeader(t, h, e)

	mustContain(t, e.Error(), "extra sound chip section")
}

func TestGetHeaderWithValidFile(t *testing.T) {
	h, e := getHeaderFromSlice(validHeader)
	hasHeaderNoError(t, h, e)
}

func TestStringWithValidPALFile(t *testing.T) {
	var palHeader = append([]byte(nil), validHeader...)
	palHeader[122] = 1 << 0

	h, e := getHeaderFromSlice(palHeader)
	hasHeaderNoError(t, h, e)

	mustContain(t, fmt.Sprint(h), ": PAL")
}

func TestStringWithValidDualRegionFile(t *testing.T) {
	var dualRegionHeader = append([]byte(nil), validHeader...)
	dualRegionHeader[122] = 1 << 1

	h, e := getHeaderFromSlice(dualRegionHeader)
	hasHeaderNoError(t, h, e)

	mustContain(t, fmt.Sprint(h), ": dual PAL/NTSC")
}

func TestStringWithMultiChipsFile(t *testing.T) {
	var multiChipsHeader = append([]byte(nil), validHeader...)
	multiChipsHeader[123] = 1<<1 | 1<<3 // VRC7 & MMC5

	h, e := getHeaderFromSlice(multiChipsHeader)
	hasHeaderNoError(t, h, e)

	mustContain(t, fmt.Sprint(h), ": VRC7, MMC5")
}

// the main() test is really just there to add a bit of code coverage
func TestMainWithValidFile(t *testing.T) {
	oldStdout := os.Stdout
	os.Stdout = nil
	defer func() { os.Stdout = oldStdout }() // not really interested in seeing main()'s stdout, so hide it for the test

	os.Args = []string{"nsfinfo", "samples/smb1.nsf"}
	main()
}

// helper method that writes a test slice to disk & then calls getHeader() on it
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

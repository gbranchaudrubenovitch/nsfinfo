package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var nsfFile = flag.String("f", "", "relative path to the nsf file you want to open. Required.")
	flag.Parse()

	if *nsfFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	header := fillHeader(*nsfFile)
	if header == nil {
		os.Exit(1) // something went wrong
	}
}

func fillHeader(relativePathToNsf string) *nsfHeader {
	f, e := os.Open(relativePathToNsf)
	defer f.Close()

	if e != nil {
		fmt.Println("nsf file cannot be found or opened.")
		return nil
	}

	// read first 128 bytes (the header only)
	rawHeader := make([]byte, 128)
	_, e = io.ReadFull(f, rawHeader)
	if e != nil {
		fmt.Println("could not read the 1st 128 bytes of the nsf file.")
		return nil
	}

	// TODO: build a nsfHeader struct from it (now just a dummy printf)
	fmt.Printf("header parsed: %q", rawHeader)
	return nil
}

type nsfHeader struct {
	// TODO: define this from nesdev wiki
}

package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("A relative path to a nsf file is required.")
	}

	var pathToFile = os.Args[1]
	header, e := getHeader(pathToFile)
	if e != nil {
		log.Fatal(e)
	}

	fmt.Print(header)
}

func getHeader(relativePathToNsf string) (*nsfHeader, error) {
	f, e := os.Open(relativePathToNsf)
	defer f.Close()

	if e != nil {
		return nil, e
	}

	// read the actual bytes
	header := new(nsfHeader)
	e = binary.Read(f, binary.BigEndian, header)
	if e != nil {
		return nil, e
	}

	// validate the header
	if string(header.Prelude[:]) != "NESM\x1A" {
		return nil, errors.New("invalid nsf file: invalid prelude")
	}

	if header.ExtraChipFlags&futureChip1 != 0 || header.ExtraChipFlags&futureChip2 != 0 {
		return nil, errors.New("invalid nsf file: extra sound chip section contains unsupported values")
	}

	return header, nil
}

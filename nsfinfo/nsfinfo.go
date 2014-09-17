package main

import (
	"encoding/binary"
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

	valid, e := header.isValid()
	if !valid {
		return nil, e
	}

	return header, nil
}

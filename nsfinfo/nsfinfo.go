package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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

const headerLen = 128

// nsfHeader exposes the nsf header fields (more info at http://wiki.nesdev.com/w/index.php/NSF)
type nsfHeader struct {
	Prelude          [5]byte
	Version          int8
	TotalSongs       int8
	StartingSong     int8
	LoadAddress      int16
	InitAddress      int16
	PlayAddress      int16
	SongName         [32]byte
	Artist           [32]byte
	CopyrightHolder  [32]byte
	NTSCPlaySpeed    int16
	BankswitchInit   [8]byte
	PalPlaySpeed     int16
	RegionFlags      regionFlag
	ExtraChipFlags   extraChipFlag
	ExpansionPadding int32 // always 0
}

func (h nsfHeader) String() string {
	return fmt.Sprintf(`	Details of %s
		song name: %s
		TODO: more fields!
		`, os.Args[1], trimNull(h.SongName[:]))
}

type regionFlag byte

const (
	pal  regionFlag = 1 << iota
	dual            // if both are off, it is a NTSC file
)

type extraChipFlag byte

const (
	vrc6 extraChipFlag = 1 << iota
	vrc7
	fds
	mmc5
	namco163
	sunsoft5b
	futureChip1 // must never be set
	futureChip2 // must never be set
)

func trimNull(s []byte) string {
	return string(s[:strings.Index(string(s), "\x00")])
}

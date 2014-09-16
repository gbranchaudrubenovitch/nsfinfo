package main

import (
	"fmt"
	"strings"
)

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
	return fmt.Sprintf("%-17s: %s\n%-17s: %s\n%-17s: %s\n%-17s: %d\n%-17s: %#x\n",
		"name", trimNull(h.SongName[:]),
		"artist", trimNull(h.Artist[:]),
		"copyright holder", trimNull(h.CopyrightHolder[:]),
		"total # of songs", h.TotalSongs,
		"load address", h.LoadAddress)
}

func trimNull(s []byte) string {
	return string(s[:strings.Index(string(s), "\x00")])
}

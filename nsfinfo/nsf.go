package main

import (
	"encoding/binary"
	"errors"
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

// nesWord is a little endian uint16
type nesWord [2]byte

func (w nesWord) toUInt16() uint16 {
	return binary.LittleEndian.Uint16(w[:])
}

// nsfHeader exposes the nsf header fields (more info at http://wiki.nesdev.com/w/index.php/NSF)
type nsfHeader struct {
	Prelude          [5]byte
	Version          int8
	TotalSongs       int8
	StartingSong     int8
	LoadAddress      nesWord
	InitAddress      nesWord
	PlayAddress      nesWord
	SongName         [32]byte
	Artist           [32]byte
	CopyrightHolder  [32]byte
	NTSCPlaySpeed    nesWord
	BankswitchInit   int64 // we do not care about the specific banks assignment
	PalPlaySpeed     nesWord
	RegionFlags      regionFlag
	ExtraChipFlags   extraChipFlag
	ExpansionPadding int32 // always 0
}

func (h nsfHeader) String() string {
	return fmt.Sprintf("%-23s: %s\n%-23s: %s\n%-23s: %s\n%-23s: %d\n%-23s: %d\n%-23s: %s\n%-23s: %d\n----------------\n%-23s: %d\n%-23s: %v\n%-23s: %s\n%-23s: %#x\n%-23s: %#x\n%-23s: %#x",
		"name", trimNull(h.SongName[:]),
		"artist", trimNull(h.Artist[:]),
		"copyright holder", trimNull(h.CopyrightHolder[:]),
		"total # of songs", h.TotalSongs,
		"first song", h.StartingSong,
		"region", h.region(),
		"play speed (μs)", h.playSpeed().toUInt16(),
		"nsf version", h.Version,
		"uses bankswitching", h.BankswitchInit != 0,
		"expansion chips in use", h.extraChips(),
		"load address", h.LoadAddress.toUInt16(),
		"init address", h.InitAddress.toUInt16(),
		"play address", h.PlayAddress.toUInt16(),
	)
}

func (h nsfHeader) playSpeed() nesWord {
	if h.RegionFlags&pal != 0 {
		return h.PalPlaySpeed
	}
	return h.NTSCPlaySpeed
}

func (h nsfHeader) region() string {
	if h.RegionFlags&dual != 0 {
		return "dual PAL/NTSC"
	} else if h.RegionFlags&pal != 0 {
		return "PAL"
	}
	return "NTSC"
}

func (h nsfHeader) extraChips() string {
	// todo: replace those ugly if/else with a const [extraChipFlag, string]map
	var chipsInUse []string
	if h.ExtraChipFlags&vrc6 != 0 {
		chipsInUse = append(chipsInUse, "VRC6")
	} else if h.ExtraChipFlags&vrc7 != 0 {
		chipsInUse = append(chipsInUse, "VRC7")
	} else if h.ExtraChipFlags&fds != 0 {
		chipsInUse = append(chipsInUse, "Famicom Disk System")
	} else if h.ExtraChipFlags&mmc5 != 0 {
		chipsInUse = append(chipsInUse, "mmc5")
	} else if h.ExtraChipFlags&namco163 != 0 {
		chipsInUse = append(chipsInUse, "Namco 163")
	} else if h.ExtraChipFlags&sunsoft5b != 0 {
		chipsInUse = append(chipsInUse, "Sunsoft 5B")
	}

	if len(chipsInUse) == 0 {
		chipsInUse = append(chipsInUse, "none")
	}

	return strings.Join(chipsInUse, ", ")
}

func (h nsfHeader) isValid() (bool, error) {
	// validate the header
	if string(h.Prelude[:]) != "NESM\x1A" {
		return false, errors.New("invalid nsf file - invalid prelude")
	}

	if h.ExtraChipFlags&futureChip1 != 0 || h.ExtraChipFlags&futureChip2 != 0 {
		return false, errors.New("invalid nsf file - extra sound chip section contains unsupported values")
	}
	return true, nil
}

func trimNull(s []byte) string {
	return string(s[:strings.Index(string(s), "\x00")])
}

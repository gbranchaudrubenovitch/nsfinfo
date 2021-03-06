package main

import (
	"errors"
	"fmt"
	"strings"
)

type regionFlag byte
type extraChipFlag byte

const (
	pal regionFlag = 1 << iota
	dual

	futureChip1 extraChipFlag = 1 << 6
	futureChip2 extraChipFlag = 1 << 7
)

var extraChips = []struct {
	flag extraChipFlag
	name string
}{
	{flag: 1 << 0, name: "VRC6"},
	{flag: 1 << 1, name: "VRC7"},
	{flag: 1 << 2, name: "Famicom Disk System"},
	{flag: 1 << 3, name: "MMC5"},
	{flag: 1 << 4, name: "Namco 163"},
	{flag: 1 << 5, name: "Sunsoft 5B"},
	{flag: futureChip1, name: "Future Chip 1 (not supported)"},
	{flag: futureChip2, name: "Future Chip 2 (not supported)"},
}

// nsfHeader exposes the nsf header fields (more info at http://wiki.nesdev.com/w/index.php/NSF)
type nsfHeader struct {
	FourCC           [5]byte
	Version          uint8
	TotalSongs       uint8
	StartingSong     uint8
	LoadAddress      uint16
	InitAddress      uint16
	PlayAddress      uint16
	SongName         [32]byte
	Artist           [32]byte
	CopyrightHolder  [32]byte
	NTSCPlaySpeed    uint16
	BankswitchInit   uint64 // we do not care about the specific banks assignment
	PalPlaySpeed     uint16
	RegionFlags      regionFlag
	ExtraChipFlags   extraChipFlag
	ExpansionPadding int32 // always 0
}

func (h nsfHeader) isValid() (bool, error) {
	if string(h.FourCC[:]) != "NESM\x1A" {
		return false, errors.New("invalid nsf file - invalid FourCC")
	}

	if h.ExtraChipFlags&futureChip1 != 0 || h.ExtraChipFlags&futureChip2 != 0 {
		return false, errors.New("invalid nsf file - extra sound chip section contains unsupported values")
	}
	return true, nil
}

func (h nsfHeader) String() string {
	return fmt.Sprintf("%-23s: %s\n%-23s: %s\n%-23s: %s\n%-23s: %d\n%-23s: %d\n%-23s: %s\n%-23s: %d\n----------------\n%-23s: %d\n%-23s: %v\n%-23s: %s\n%-23s: %#X\n%-23s: %#X\n%-23s: %#X",
		"name", trimNull(h.SongName[:]),
		"artist", trimNull(h.Artist[:]),
		"copyright holder", trimNull(h.CopyrightHolder[:]),
		"total # of songs", h.TotalSongs,
		"first song", h.StartingSong,
		"region", h.region(),
		"play speed (μs)", h.playSpeed(),
		"nsf version", h.Version,
		"uses bankswitching", h.BankswitchInit != 0,
		"expansion chips in use", h.extraChips(),
		"load address", h.LoadAddress,
		"init address", h.InitAddress,
		"play address", h.PlayAddress,
	)
}

func (h nsfHeader) playSpeed() uint16 {
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
	var chipsInUse []string
	for _, i := range extraChips {
		if h.ExtraChipFlags&i.flag != 0 {
			chipsInUse = append(chipsInUse, i.name)
		}
	}

	if len(chipsInUse) == 0 {
		chipsInUse = append(chipsInUse, "none")
	}

	return strings.Join(chipsInUse, ", ")
}

func trimNull(s []byte) string {
	return string(s[:strings.Index(string(s), "\x00")])
}

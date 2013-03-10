package main

import (
	"os"
)

type Rom struct {
	PrgRomCount int
	ChrRomCount int
	PrgRoms     [][]byte
	ChrRoms     [][]byte
}

func LoadRom(filename string) *Rom {
	rom := Rom{}
	file, err := os.Open(filename)
	if err != nil {
		fatal("file not found: ", filename)
	}

	header := make([]byte, 16)
	count, err := file.Read(header)
	if err != nil || count != 16 {
		fatal("can't read header")
	}

	debug("header %p\n", header)

	if string(header[0:3]) != "NES" {
		fatal("invalid file format")
	}

	rom.PrgRomCount = int(header[4])
	debug("prg count %d\n", rom.PrgRomCount)
	rom.ChrRomCount = int(header[5])
	debug("chr count %d\n", rom.ChrRomCount)

	mapper := (header[6]&0xF0)>>4 | (header[7] & 0xF0)

	mirroring := ""

	if header[6]&0x01 == 1 {
		mirroring = "vertical"
	} else {
		mirroring = "horizontal"
	}

	debug("mirroring %d", mirroring)

	debug("mapper %d", mapper)

	debug("control")

	rom.PrgRoms = make([][]byte, rom.PrgRomCount)
	for i := 0; i < rom.PrgRomCount; i++ {
		rom.PrgRoms[i] = make([]byte, 0x4000)
		count, err := file.Read(rom.PrgRoms[i])
		if err != nil || count != 0x4000 {
			fatal("can't read prg rom:", i+1)
		}
	}

	rom.ChrRoms = make([][]byte, rom.ChrRomCount)
	for i := 0; i < rom.ChrRomCount; i++ {
		rom.ChrRoms[i] = make([]byte, 0x2000)
		count, err := file.Read(rom.ChrRoms[i])
		if err != nil || count != 0x2000 {
			fatal("can't read chr rom:", i+1)
		}
	}

	return &rom
}

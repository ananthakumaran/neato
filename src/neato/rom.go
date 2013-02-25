package main

import (
	"fmt"
	"os"
)

type Rom struct {
	PrgRomCount int
	ChrRomCount int
	PrgRoms     [][]byte
	ChrRoms     [][]byte
}

func LoadRom(filename string) Rom {
	file, err := os.Open(filename)
	if err != nil {
		fatal("file not found: ", filename)
	}

	header := make([]byte, 16)
	count, err := file.Read(header)
	if err != nil || count != 16 {
		fatal("can't read header")
	}

	fmt.Println(header)

	if string(header[0:3]) != "NES" {
		fatal("invalid file format")
	}

	PrgRomCount := header[4]
	fmt.Println("prg count", PrgRomCount)
	ChrRomCount := header[5]
	fmt.Println("chr count", ChrRomCount)

	mapper := (header[6]&0xF0)>>4 | (header[7] & 0xF0)

	fmt.Println("mapper ", mapper)

	fmt.Println("control")

	rom := Rom{}

	rom.PrgRoms = make([][]byte, PrgRomCount)
	for i := byte(0); i < PrgRomCount; i++ {
		rom.PrgRoms[i] = make([]byte, 0x4000)
		count, err := file.Read(rom.PrgRoms[i])
		if err != nil || count != 0x4000 {
			fatal("can't read prg rom:", i+1)
		}
	}

	rom.ChrRoms = make([][]byte, ChrRomCount)
	for i := byte(0); i < ChrRomCount; i++ {
		rom.ChrRoms[i] = make([]byte, 0x2000)
		count, err := file.Read(rom.ChrRoms[i])
		if err != nil || count != 0x2000 {
			fatal("can't read chr rom:", i+1)
		}
	}

	return rom
}
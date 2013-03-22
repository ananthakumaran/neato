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
	mirroring   int
}

type Mapper interface {
	PrgRead(uint16) byte
	PrgWrite(uint16, byte)
	ChrRead(uint16) byte
	ChrWrite(uint16, byte)
	Mirroring() int
}

type NROM struct {
	rom *Rom
	prg []byte
	chr []byte
}

func newNROM(rom *Rom) *NROM {
	nrom := &NROM{rom: rom}
	nrom.prg = make([]byte, 0x8000)
	copy(nrom.prg[0x0:0x4000], rom.PrgRoms[0])

	switch rom.PrgRomCount {
	case 1:
		copy(nrom.prg[0x4000:0x8000], rom.PrgRoms[0])
	case 2:
		copy(nrom.prg[0x4000:0x8000], rom.PrgRoms[1])
	default:
		fatal("uknown prg rom count", rom.PrgRomCount)
	}

	nrom.chr = make([]byte, 0x2000)
	switch rom.ChrRomCount {
	case 0: // no chr
	case 1:
		copy(nrom.chr, rom.ChrRoms[0])
	}

	return nrom
}

func (nrom *NROM) PrgRead(address uint16) byte {
	return nrom.prg[address]
}

func (nrom *NROM) PrgWrite(address uint16, val byte) {
}

func (nrom *NROM) ChrRead(address uint16) byte {
	return nrom.chr[address]
}

func (nrom *NROM) ChrWrite(address uint16, val byte) {
	nrom.chr[address] = val
}

func (nrom *NROM) Mirroring() int {
	return nrom.rom.mirroring
}

type UxROM struct {
	rom         *Rom
	chr         []byte
	currentBank uint8
}

func newUxROM(rom *Rom) *UxROM {
	uxrom := &UxROM{rom: rom}

	uxrom.chr = make([]byte, 0x2000)
	switch rom.ChrRomCount {
	case 0: // no chr
	case 1:
		copy(uxrom.chr, rom.ChrRoms[0])
	}

	return uxrom
}

func (uxrom *UxROM) PrgRead(address uint16) byte {
	if address >= 0x4000 {
		return uxrom.rom.PrgRoms[uxrom.rom.PrgRomCount-1][address-0x4000]
	}
	return uxrom.rom.PrgRoms[uxrom.currentBank][address]
}

func (uxrom *UxROM) PrgWrite(address uint16, val byte) {
	uxrom.currentBank = uint8(val) & uint8(uxrom.rom.PrgRomCount-1)
}

func (uxrom *UxROM) ChrRead(address uint16) byte {
	return uxrom.chr[address]
}

func (uxrom *UxROM) ChrWrite(address uint16, val byte) {
	uxrom.chr[address] = val
}

func (uxrom *UxROM) Mirroring() int {
	return uxrom.rom.mirroring
}

func LoadRom(filename string) Mapper {
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

	if string(header[0:3]) != "NES" {
		fatal("invalid file format")
	}

	if string(header[7:0xF]) == "DiskDude" ||
		string(header[7:0xF]) == "demiforce" {
		copy(header[7:0xF], make([]byte, 8))
	}

	debug("header %v\n", header)

	rom.PrgRomCount = int(header[4])
	rom.ChrRomCount = int(header[5])

	mapper := header[6]>>4 | (header[7] & 0xF0)

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

	fmt.Printf("PRG ROM: %d x 16KiB\n", rom.PrgRomCount)
	fmt.Printf("CHR ROM: %d x  8KiB\n", rom.ChrRomCount)
	fmt.Printf("Mapper : #%d\n", mapper)

	if header[6]&0x01 == 1 {
		rom.mirroring = VERTICAL_MIRRORING
		fmt.Println("Mirroring: Vertical")
	} else {
		rom.mirroring = HORIZONTAL_MIRRORING
		fmt.Println("Mirroring: Horizontal")
	}

	switch mapper {
	case 0:
		return newNROM(&rom)
	case 2:
		return newUxROM(&rom)
	default:
		fatal("unimplemented mapper", mapper)
	}
	return nil
}

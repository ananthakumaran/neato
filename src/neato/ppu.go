package main

type Ppu struct {
	rom    Rom
	ram    []byte
	status uint8
	// 7 6 5 4
	// | | | -- ignore write to vram
	// | | -- more than 8 sprite on the current scanline
	// | -- sprite hit flag
	// -- V-Blank is occurring

	// control registers
	nameTableAddress      uint16
	incrementBy           uint8
	spritePatternTable    uint16
	backgroundPatterTable uint16
	spriteSize            int
	nmiOnVBlank           bool
	color                 bool
	monochrome            bool
	showclipBackground    bool
	showclipSprite        bool
	displayBackground     bool
	displaySprite         bool
	colorIntensity        uint8
}

func newPpu(rom Rom) Ppu {
	ppu := Ppu{}
	ppu.rom = rom
	ppu.ram = make([]byte, 0x10000)
	ppu.vBlank(true)
	return ppu
}

func (ppu *Ppu) vBlank(set bool) {
	if set {
		ppu.status |= 0x80
	} else {
		ppu.status |= 0x7F
	}
}

// 7 6 5 4 3 2 1 0
// | | | | | | --|  name table address 
// | | | | | | ---  0 -> $2000, 1 -> $2400, 2 -> $2800, 3 -> $2C00
// | | | | | -- amount to increment 0 -> 1, 1 -> 32
// | | | | -- sprite pattern table 0 -> $0000, 1 -> $1000
// | | | -- background pattern table 0 -> $0000, 1 -> $1000
// | | -- size of sprite in pixels 0 -> 8x8, 1 -> 8x16
// | -- unused
// -- nmi on V-Blank
func (ppu *Ppu) controlRegister1(val uint8) {

	switch val & 0x03 {
	case 0:
		ppu.nameTableAddress = 0x2000
	case 1:
		ppu.nameTableAddress = 0x2400
	case 2:
		ppu.nameTableAddress = 0x2800
	case 3:
		ppu.nameTableAddress = 0x2C00
	}

	if (val>>2)&1 == 0 {
		ppu.incrementBy = 1
	} else {
		ppu.incrementBy = 32
	}

	if (val>>3)&1 == 0 {
		ppu.spritePatternTable = 0x0000
	} else {
		ppu.spritePatternTable = 0x1000
	}

	if (val>>4)&1 == 0 {
		ppu.backgroundPatterTable = 0x0000
	} else {
		ppu.backgroundPatterTable = 0x1000
	}

	if (val>>5)&1 == 0 {
		ppu.spriteSize = 8
	} else {
		ppu.spriteSize = 16
	}

	ppu.nmiOnVBlank = (val>>7 == 1)
}

// 7 6 5 4 3 2 1 0
// | | | | | | | -- 0 -> color, 1 -> monochrome
// | | | | | | ---  clip the background
// | | | | | -- clip sprites
// | | | | -- show background
// | | | -- background pattern table 0 -> $0000, 1 -> $1000
// -------- background color in monochrome or color intensity in color mode

func (ppu *Ppu) controlRegister2(val uint8) {
	if (val & 1) == 0 {
		ppu.color = true
	} else {
		ppu.monochrome = true
	}

	ppu.showclipBackground = (val>>1)&1 == 1
	ppu.showclipSprite = (val>>2)&1 == 1

	ppu.displayBackground = (val>>3)&1 == 1
	ppu.displaySprite = (val>>4)&1 == 1
	ppu.colorIntensity = val >> 5
}

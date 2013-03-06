package main

type Ppu struct {
	cpu    *Cpu
	rom    *Rom
	ram    *Memory
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

	// vram
	addrStatus int
	address    uint16

	scrollAddrStatus int
	scrollAddress    uint16

	// oram
	oamAddress uint16
	oamRam     *Memory
}

func newPpu(rom *Rom) *Ppu {
	ppu := Ppu{}
	ppu.rom = rom
	ppu.ram = newMemory(0xFFFF)
	ppu.ram.copy(0x0000, 0x2000, rom.ChrRoms[0])
	ppu.ram.mirror(0x2000, 0x2EFF, 0x3000, 0x3EFF)
	ppu.ram.mirror(0x3F00, 0x3F1F, 0x3F20, 0x3FFF)
	ppu.ram.mirror(0x0000, 0x3FFF, 0x4000, 0xFFFF)
	ppu.oamRam = newMemory(0xFF)
	ppu.vBlank(true)
	ppu.reset()
	return &ppu
}

func (ppu *Ppu) reset() {
	ppu.addrStatus = 0
	ppu.scrollAddrStatus = 0
	ppu.controlRegister1(0)
	ppu.controlRegister2(0)
}

func (ppu *Ppu) vBlank(set bool) {
	if set {
		ppu.status |= 0x80
	} else {
		ppu.status &= 0x7F
	}
}

func (ppu *Ppu) read(address uint16) byte {
	debug("R %X\n", address)

	switch address {
	case 0x2002:
		status := ppu.status
		ppu.vBlank(false)
		ppu.resetLatch()
		debug("status %b\n", status)
		return status

	case 0x2004:
	case 0x2007:
		address := ppu.address
		ppu.address += uint16(ppu.incrementBy)
		return ppu.ram.read(address)
	default:
		fatal("unimplemented")
	}

	return 0
}

func (ppu *Ppu) write(address uint16, val byte) {
	debug("W %X %X \n", address, val)

	switch address {
	case 0x2000:
		ppu.controlRegister1(val)
	case 0x2001:
		ppu.controlRegister2(val)
	case 0x2003:
		ppu.oamAddress = uint16(val)
	case 0x2004:
		debug("W %X %X \n", address, val)
		ppu.oamRam.write(uint16(ppu.oamAddress), val)
		ppu.oamAddress++
	case 0x2005:
		debug("W %X %X \n", address, val)
		switch ppu.scrollAddrStatus {
		case 0:
			ppu.scrollAddress = (8 << uint16(val)) | ppu.scrollAddress&0x00FF
			ppu.scrollAddrStatus++
		case 1:
			ppu.scrollAddress = uint16(val) | ppu.scrollAddress&0xFF00
			ppu.scrollAddrStatus = 0
		}
	case 0x2006:
		switch ppu.addrStatus {
		case 0:
			ppu.address = (uint16(val) << 8) | ppu.address&0x00FF
			//debug(" VRAM %X val %X  ", ppu.address, val)
			ppu.addrStatus++
		case 1:
			ppu.address = uint16(val) | ppu.address&0xFF00
			//debug(" VRAM %X val %X ", ppu.address, val)
			ppu.addrStatus = 0
		}
	case 0x2007:
		if ppu.address > 0x3FFF {
			fatal("mirror not implemented")
		}

		//debug(" VRAM %X val %X  ", ppu.address, val)
		ppu.ram.write(ppu.address, val)
		ppu.address += uint16(ppu.incrementBy)
	case 0x4014:
		debug("\n OAM DMC \n")
		base := uint16(val) * 0x100
		for i := ppu.oamAddress; i <= 255; i++ {
			ppu.oamRam.write(i, ppu.cpu.ram.read(base+uint16(i)))
		}

	case 0x4011, 0x4015:
	default:
		fatal("unimplemented")
	}
}

func (ppu *Ppu) resetLatch() {
	ppu.addrStatus = 0
	ppu.scrollAddrStatus = 0
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
	debug("PPU CTRL %x ", val)

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
// | | | -- show sprite
// -------- background color in monochrome or color intensity in color mode

func (ppu *Ppu) controlRegister2(val uint8) {
	debug("PPU MASK %x ", val)
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

func (ppu *Ppu) drawScreen() {
	debug("draw screen\n")
	debug("name table addres %x \n", ppu.nameTableAddress)

	debug("\nsprite table \n")
	for i := 0x0; i <= 0x0FFF; {
		for j := 0; j < 32; j++ {
			debug("%x ", ppu.ram.read(uint16(i)))
			i++
		}
		debug("\n")
	}

	debug("background table \n")
	for i := 0x1000; i <= 0x1FFF; {
		for j := 0; j < 32; j++ {
			debug("%x ", ppu.ram.read(uint16(i)))
			i++
		}
		debug("\n")
	}

	debug("name and table attributes \n")
	for i := 0x2000; i <= 0x2FFF; {
		for j := 0; j < 32; j++ {
			debug("%x ", ppu.ram.read(uint16(i)))
			i++
		}
		debug("\n")
	}

	debug("ooam\n")
	for i := 0; i < 255; {
		for j := 0; j < 32; j++ {
			debug("%x ", ppu.oamRam.read(uint16(i)))
			i++
		}
		debug("\n")
	}
}

func (ppu *Ppu) startVblank() {
	ppu.drawScreen()
	ppu.vBlank(true)
	debug("\n nmi request \n")
	ppu.cpu.Interrupt(NMI)
}

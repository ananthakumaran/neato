package main

const (
	TILES_PER_ROW              = 32
	TILES_PER_COLUMN           = 30
	HORIZONTAL_PIXELS_PER_TILE = 8
	VERTICAL_PIXES_PER_TILE    = 8

	PATTERN_BYTES_PER_TILE = 16
)

var colorPalette = [][]byte{
	{0x75, 0x75, 0x75},
	{0x27, 0x1B, 0x8F},
	{0x00, 0x00, 0xAB},
	{0x47, 0x00, 0x9F},
	{0x8F, 0x00, 0x77},
	{0xAB, 0x00, 0x13},
	{0xA7, 0x00, 0x00},
	{0x7F, 0x0B, 0x00},
	{0x43, 0x2F, 0x00},
	{0x00, 0x47, 0x00},
	{0x00, 0x51, 0x00},
	{0x00, 0x3F, 0x17},
	{0x1B, 0x3F, 0x5F},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0xBC, 0xBC, 0xBC},
	{0x00, 0x73, 0xEF},
	{0x23, 0x3B, 0xEF},
	{0x83, 0x00, 0xF3},
	{0xBF, 0x00, 0xBF},
	{0xE7, 0x00, 0x5B},
	{0xDB, 0x2B, 0x00},
	{0xCB, 0x4F, 0x0F},
	{0x8B, 0x73, 0x00},
	{0x00, 0x97, 0x00},
	{0x00, 0xAB, 0x00},
	{0x00, 0x93, 0x3B},
	{0x00, 0x83, 0x8B},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0xFF, 0xFF, 0xFF},
	{0x3F, 0xBF, 0xFF},
	{0x5F, 0x97, 0xFF},
	{0xA7, 0x8B, 0xFD},
	{0xF7, 0x7B, 0xFF},
	{0xFF, 0x77, 0xB7},
	{0xFF, 0x77, 0x63},
	{0xFF, 0x9B, 0x3B},
	{0xF3, 0xBF, 0x3F},
	{0x83, 0xD3, 0x13},
	{0x4F, 0xDF, 0x4B},
	{0x58, 0xF8, 0x98},
	{0x00, 0xEB, 0xDB},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0xFF, 0xFF, 0xFF},
	{0xAB, 0xE7, 0xFF},
	{0xC7, 0xD7, 0xFF},
	{0xD7, 0xCB, 0xFF},
	{0xFF, 0xC7, 0xFF},
	{0xFF, 0xC7, 0xDB},
	{0xFF, 0xBF, 0xB3},
	{0xFF, 0xDB, 0xAB},
	{0xFF, 0xE7, 0xA3},
	{0xE3, 0xFF, 0xA3},
	{0xAB, 0xF3, 0xBF},
	{0xB3, 0xFF, 0xCF},
	{0x9F, 0xFF, 0xF3},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00}}

type Ppu struct {
	cpu  *Cpu
	rom  *Rom
	vram *Memory

	// CRTL
	ctrlRegister                 uint8
	basenameTableAddress         uint16
	incrementBy                  uint8
	spritePatternTableAddress    uint16
	backgroundPatterTableAddress uint16
	spriteSize                   int
	nmiOnVBlank                  bool

	// MASK
	maskRegister       uint8
	color              bool
	monochrome         bool
	showclipBackground bool
	showclipSprite     bool
	displayBackground  bool
	displaySprite      bool
	colorIntensity     uint8

	// STATUS
	status          uint8
	fSpriteOverflow bool
	fSpritZeroHit   bool
	fVerticalBlank  bool

	// vram
	addrStatus int
	address    uint16

	scrollAddrStatus int
	scrollAddress    uint16

	// oram
	oamAddress uint16
	oamRam     *Memory

	// bookeeping
	scanline             int
	currentScanlineCycle int
}

func newPpu(rom *Rom) *Ppu {
	ppu := Ppu{}
	ppu.rom = rom
	ppu.vram = newMemory(0xFFFF)

	if rom.ChrRomCount == 1 {
		ppu.vram.copy(0x0000, 0x2000, rom.ChrRoms[0])
	}

	ppu.vram.mirror(0x2000, 0x2EFF, 0x3000, 0x3EFF)
	ppu.vram.mirror(0x3F00, 0x3F1F, 0x3F20, 0x3FFF)
	ppu.vram.mirror(0x0000, 0x3FFF, 0x4000, 0xFFFF)
	ppu.oamRam = newMemory(0xFF)
	//ppu.vBlank(true)
	ppu.reset()
	return &ppu
}

func (ppu *Ppu) reset() {
	ppu.addrStatus = 0
	ppu.scrollAddrStatus = 0
	ppu.controlRegister1(0)
	ppu.controlRegister2(0)
	ppu.oamAddress = 0

	ppu.currentScanlineCycle = 0
	ppu.scanline = 241
}

func (ppu *Ppu) read(address uint16) byte {
	debug("R %X\n", address)

	switch address {
	case 0x2000:
		return ppu.ctrlRegister
	case 0x2001:
		return ppu.maskRegister
	case 0x2002:
		status := ppu.getStatus()
		ppu.fVerticalBlank = false
		ppu.resetLatch()
		debug("status %b\n", status)
		return status
	case 0x2004:
	case 0x2007:
		address := ppu.address
		ppu.address += uint16(ppu.incrementBy)
		return ppu.vram.read(address)
		// joystick
	case 0x4016, 0x4017:
		return 0
	default:
		info("READ unimplemented %X", address)
		//fatal("READ unimplemented %X", address)
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
	case 0x2002:
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
		debug(" VRAM %X val %X  ", ppu.address, val)
		ppu.vram.write(ppu.address, val)
		ppu.address += uint16(ppu.incrementBy)
	case 0x4014:
		debug("\n OAM DMC \n")
		base := uint16(val) * 0x100
		for i := ppu.oamAddress; i <= 255; i++ {
			ppu.oamRam.write(i, ppu.cpu.ram.read(base+uint16(i)))
		}

	case 0x4015, 0x4017:
	case 0x4000, 0x4001, 0x4002, 0x4003, 0x4004,
		0x4005, 0x4006, 0x4007, 0x4008, 0x4009, 0x400A,
		0x400B, 0x400C, 0x400D, 0x400E, 0x400F, 0x4010,
		0x4011, 0x4012, 0x4013:
		info("sound beep")
	default:
		info("write unimplemented %x %x", address, val)
		//fatal("write unimplemented %x %x", address, val)
	}
}

func (ppu *Ppu) resetLatch() {
	ppu.addrStatus = 0
	ppu.scrollAddrStatus = 0
}

// 76543210
// ||||||||
// |||+++++- Least significant bits previously written into a PPU register
// |||       (due to register not being updated for this address)
// ||+------ Sprite overflow. The PPU can handle only eight sprites on one
// ||        scanline and sets this bit if it starts dropping sprites.
// ||        Normally, this triggers when there are 9 sprites on a scanline,
// ||        but the actual behavior is significantly more complicated.
// |+------- Sprite 0 Hit.  Set when a nonzero pixel of sprite 0 overlaps
// |         a nonzero background pixel, cleared at start of pre-render line.
// |         Used for raster timing.
// +-------- Vertical blank has started (0: not in VBLANK; 1: in VBLANK)

func (ppu *Ppu) getStatus() uint8 {
	ppu.status &= 0x1F
	if ppu.fSpriteOverflow {
		ppu.status |= 0x20
	}

	if ppu.fSpritZeroHit {
		ppu.status |= 0x40
	}

	if ppu.fVerticalBlank {
		ppu.status |= 0x80
	}

	return ppu.status
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
	debug("PPU CTRL %02X ", val)
	ppu.ctrlRegister = val

	switch val & 0x03 {
	case 0:
		ppu.basenameTableAddress = 0x2000
	case 1:
		ppu.basenameTableAddress = 0x2400
	case 2:
		ppu.basenameTableAddress = 0x2800
	case 3:
		ppu.basenameTableAddress = 0x2C00
	}

	if (val>>2)&1 == 0 {
		ppu.incrementBy = 1
	} else {
		ppu.incrementBy = 32
	}

	if (val>>3)&1 == 0 {
		ppu.spritePatternTableAddress = 0x0000
	} else {
		ppu.spritePatternTableAddress = 0x1000
	}

	if (val>>4)&1 == 0 {
		ppu.backgroundPatterTableAddress = 0x0000
	} else {
		ppu.backgroundPatterTableAddress = 0x1000
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
	ppu.maskRegister = val
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
	debug("name table addres %x \n", ppu.basenameTableAddress)

	debug("\nsprite table \n")
	for i := 0x0; i <= 0x0FFF; {
		for j := 0; j < 32; j++ {
			debug("%02X ", ppu.vram.read(uint16(i)))
			i++
		}
		debug("\n")
	}

	debug("background table \n")
	for i := 0x1000; i <= 0x1FFF; {
		for j := 0; j < 32; j++ {
			debug("%02X ", ppu.vram.read(uint16(i)))
			i++
		}
		debug("\n")
	}

	debug("name and table attributes \n")
	for i := 0x2000; i <= 0x2FFF; {
		for j := 0; j < 32; j++ {
			debug("%02X ", ppu.vram.read(uint16(i)))
			i++
		}
		debug("\n")
	}

	debug("ooam\n")
	for i := 0; i < 255; {
		for j := 0; j < 32; j++ {
			debug("%02X ", ppu.oamRam.read(uint16(i)))
			i++
		}
		debug("\n")
	}
}

func (ppu *Ppu) startVblank() {
	//	ppu.drawScreen()
	ppu.fVerticalBlank = true
	ppu.cpu.Interrupt(NMI)
}

func (ppu *Ppu) patternColor(x, y int) uint8 {
	nametableOffset := uint16((x / HORIZONTAL_PIXELS_PER_TILE) +
		((y / VERTICAL_PIXES_PER_TILE) * TILES_PER_ROW))

	patternTileNumber := ppu.vram.read(ppu.basenameTableAddress + nametableOffset)

	patternAddress := ppu.backgroundPatterTableAddress + uint16((int(patternTileNumber)*PATTERN_BYTES_PER_TILE)+(y%VERTICAL_PIXES_PER_TILE))

	pattern1 := ppu.vram.read(patternAddress)
	pattern2 := ppu.vram.read(patternAddress + 8)

	bitOffset := uint8((x % 8))
	pattern1 = (pattern1 << bitOffset) >> 7
	pattern2 = (((pattern2 << bitOffset) >> 7) << 1)
	patternIndex := pattern1 + pattern2
	return patternIndex
}

var attributeTableLookup = [][]uint8{
	{0x0, 0x1, 0x4, 0x5},
	{0x2, 0x3, 0x6, 0x7},
	{0x8, 0x9, 0xC, 0xD},
	{0xA, 0xB, 0xE, 0xF}}

func (ppu *Ppu) attributeColor(x, y int) uint8 {
	tileNumber := uint16((x / HORIZONTAL_PIXELS_PER_TILE) +
		((y / VERTICAL_PIXES_PER_TILE) * TILES_PER_ROW))

	tileRowNumber := tileNumber / TILES_PER_ROW
	tileColumnNumber := tileNumber % TILES_PER_ROW
	horizontalOffset := tileColumnNumber / 4
	verticalOffset := (tileRowNumber / 4) * 8

	attributeByte := ppu.vram.read(ppu.basenameTableAddress + 960 + horizontalOffset + verticalOffset)

	horizontalOffset = tileNumber % 4
	verticalOffset = tileRowNumber % 4

	attributeTileNumber := attributeTableLookup[verticalOffset][horizontalOffset]
	attributeColorIndex := uint8(0)

	switch attributeTileNumber {
	case 0x0, 0x1, 0x2, 0x3:
		attributeColorIndex = attributeByte << 6
	case 0x4, 0x5, 0x6, 0x7:
		attributeColorIndex = attributeByte << 4
	case 0x8, 0x9, 0xA, 0xB:
		attributeColorIndex = attributeByte << 2
	case 0xC, 0xD, 0xE, 0xF:
		attributeColorIndex = attributeByte
	}

	return (attributeColorIndex >> 6)
}

func (ppu *Ppu) renderPixel() {
	//debug("basetable %04X\n", ppu.basenameTableAddress)
	x := ppu.currentScanlineCycle
	y := ppu.scanline

	patternIndex := ppu.patternColor(x, y)
	attributeIndex := ppu.attributeColor(x, y)

	backgroundColourIndex := uint8(0)

	if patternIndex != 0 {
		backgroundColourIndex = patternIndex + (attributeIndex << 2)
	}

	if !ppu.displayBackground {
		backgroundColourIndex = 0
	}

	colorIndex := ppu.vram.read(0x3F00 + uint16(backgroundColourIndex))
	color := colorPalette[colorIndex]
	//debug("%x", color[0])
	DrawPixel(x, y, color[0], color[1], color[2])
	//debug("x %d y %x rgb # %x%x%x\n", x, y, color[0], color[1], color[2])
}

func (ppu *Ppu) step() {
	ppu.currentScanlineCycle++
	if ppu.currentScanlineCycle == 341 {
		//		debug("scan line %x\n", ppu.scanline)
		ppu.scanline++
		if ppu.scanline == 261 {
			ppu.scanline = -1
			RefreshScreen()
		}
		ppu.currentScanlineCycle = 0

		if ppu.scanline == 241 {
			ppu.startVblank()
		}

		if ppu.scanline == 0 {
			ppu.fVerticalBlank = false
			ppu.fSpriteOverflow = false
			ppu.fSpritZeroHit = false
		}
	}

	if ppu.scanline >= 0 && ppu.scanline <= 239 {
		ppu.renderPixel()
	}
}

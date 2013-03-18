package main

const (
	TILES_PER_ROW              = 32
	TILES_PER_COLUMN           = 30
	HORIZONTAL_PIXELS_PER_TILE = 8
	VERTICAL_PIXES_PER_TILE    = 8

	PATTERN_BYTES_PER_TILE = 16

	SCREEN_WIDTH  = 256
	SCREEN_HEIGHT = 240
)

const (
	VERTICAL_MIRRORING = iota
	HORIZONTAL_MIRRORING
)

// http://www.thealmightyguru.com/Games/Hacking/Wiki/index.php?title=NES_Palette
var colorPalette = [][]byte{
	{0x7C, 0x7C, 0x7C},
	{0x00, 0x00, 0xFC},
	{0x00, 0x00, 0xBC},
	{0x44, 0x28, 0xBC},
	{0x94, 0x00, 0x84},
	{0xA8, 0x00, 0x20},
	{0xA8, 0x10, 0x00},
	{0x88, 0x14, 0x00},
	{0x50, 0x30, 0x00},
	{0x00, 0x78, 0x00},
	{0x00, 0x68, 0x00},
	{0x00, 0x58, 0x00},
	{0x00, 0x40, 0x58},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0xBC, 0xBC, 0xBC},
	{0x00, 0x78, 0xF8},
	{0x00, 0x58, 0xF8},
	{0x68, 0x44, 0xFC},
	{0xD8, 0x00, 0xCC},
	{0xE4, 0x00, 0x58},
	{0xF8, 0x38, 0x00},
	{0xE4, 0x5C, 0x10},
	{0xAC, 0x7C, 0x00},
	{0x00, 0xB8, 0x00},
	{0x00, 0xA8, 0x00},
	{0x00, 0xA8, 0x44},
	{0x00, 0x88, 0x88},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0xF8, 0xF8, 0xF8},
	{0x3C, 0xBC, 0xFC},
	{0x68, 0x88, 0xFC},
	{0x98, 0x78, 0xF8},
	{0xF8, 0x78, 0xF8},
	{0xF8, 0x58, 0x98},
	{0xF8, 0x78, 0x58},
	{0xFC, 0xA0, 0x44},
	{0xF8, 0xB8, 0x00},
	{0xB8, 0xF8, 0x18},
	{0x58, 0xD8, 0x54},
	{0x58, 0xF8, 0x98},
	{0x00, 0xE8, 0xD8},
	{0x78, 0x78, 0x78},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00},
	{0xFC, 0xFC, 0xFC},
	{0xA4, 0xE4, 0xFC},
	{0xB8, 0xB8, 0xF8},
	{0xD8, 0xB8, 0xF8},
	{0xF8, 0xB8, 0xF8},
	{0xF8, 0xA4, 0xC0},
	{0xF0, 0xD0, 0xB0},
	{0xFC, 0xE0, 0xA8},
	{0xF8, 0xD8, 0x78},
	{0xD8, 0xF8, 0x78},
	{0xB8, 0xF8, 0xB8},
	{0xB8, 0xF8, 0xD8},
	{0x00, 0xFC, 0xFC},
	{0xF8, 0xD8, 0xF8},
	{0x00, 0x00, 0x00},
	{0x00, 0x00, 0x00}}

type Ppu struct {
	cpu  *Cpu
	rom  *Rom
	vram *Memory
	gui  *Gui

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
	scrollX          uint8
	scrollY          uint8

	// oram
	oamAddress uint16
	oamRam     *Memory

	// bookeeping
	scanline             int
	currentScanlineCycle int
	sprites              []Sprite

	vramReadBuffer uint8

	x          uint16
	fineX      uint16
	y          uint16
	scrollBase uint16
}

type Sprite struct {
	visible          bool
	behindBackground bool
	paletteIndex     uint8
}

func newPpu(rom *Rom) *Ppu {
	ppu := Ppu{}
	ppu.gui = newGui()
	ppu.rom = rom
	ppu.vram = newMemory(0xFFFF)

	switch rom.ChrRomCount {
	case 0: // no chr rom
	case 1:
		ppu.vram.copy(0x0000, 0x2000, rom.ChrRoms[0])
	default:
		fatal("unknown chrom count", rom.ChrRomCount)
	}

	ppu.vram.mirror(0x2000, 0x2EFF, 0x3000, 0x3EFF)
	ppu.vram.mirror(0x3F00, 0x3F1F, 0x3F20, 0x3FFF)
	ppu.vram.mirror(0x0000, 0x3FFF, 0x4000, 0xFFFF)

	switch rom.mirroring {
	case VERTICAL_MIRRORING:
		ppu.vram.mirror(0x2000, 0x23FF, 0x2800, 0x2BFF)
		ppu.vram.mirror(0x2400, 0x27FF, 0x2C00, 0x2FFF)
	case HORIZONTAL_MIRRORING:
		ppu.vram.mirror(0x2000, 0x23FF, 0x2400, 0x27FF)
		ppu.vram.mirror(0x2800, 0x2BFF, 0x2C00, 0x2FFF)
	}

	// pallete mirroring
	ppu.vram.mirror(0x3F10, 0x3F10, 0x3F00, 0x3F00)
	ppu.vram.mirror(0x3F14, 0x3F14, 0x3F04, 0x3F04)
	ppu.vram.mirror(0x3F18, 0x3F18, 0x3F08, 0x3F08)
	ppu.vram.mirror(0x3F1C, 0x3F1C, 0x3F0C, 0x3F0C)

	ppu.oamRam = newMemory(0xFF)
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
		return status
	case 0x2004:
		return ppu.oamRam.read(uint16(ppu.oamAddress))
	case 0x2007:
		address := ppu.address
		ppu.address += uint16(ppu.incrementBy)
		buffered := ppu.vramReadBuffer
		if address >= 0x3F00 && address <= 0x3FFF {
			ppu.vramReadBuffer = ppu.vram.read(address - 0x1000)
			return ppu.vram.read(address)
		}
		ppu.vramReadBuffer = ppu.vram.read(address)
		return buffered
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
			ppu.scrollX = val
			ppu.fineX = uint16(val & 0x7)
			ppu.x &= 0xFFF8
			ppu.x += ppu.fineX
			ppu.scrollAddrStatus++
		case 1:
			ppu.scrollY = val
			ppu.scrollAddrStatus = 0
		}
	case 0x2006:
		switch ppu.addrStatus {
		case 0:
			ppu.address = (uint16(val) << 8) | ppu.address&0x00FF
			ppu.addrStatus++
		case 1:
			ppu.address = uint16(val) | ppu.address&0xFF00
			ppu.addrStatus = 0

			ppu.x = ((ppu.address & 0x001F) << 3) + ppu.fineX
			ppu.scrollX = uint8(ppu.x)
			ppu.y = (((ppu.address >> 5) & 0x001F) << 3) + ((ppu.address >> 12) & 0x7)
			ppu.scrollY = uint8(ppu.y)
			switch (ppu.address >> 10) & 0x03 {
			case 0:
				ppu.scrollBase = 0x2000
			case 1:
				ppu.scrollBase = 0x2400
			case 2:
				ppu.scrollBase = 0x2800
			case 3:
				ppu.scrollBase = 0x2C00
			}
			ppu.basenameTableAddress = ppu.scrollBase
		}
	case 0x2007:
		debug(" VRAM %X val %X  ", ppu.address, val)
		ppu.vram.write(ppu.address, val)
		ppu.address += uint16(ppu.incrementBy)
	case 0x4014:
		debug("\n OAM DMC \n")
		base := uint16(val) * 0x100
		addr := uint8(ppu.oamAddress)

		for i := 0; i <= 255; i++ {
			ppu.oamRam.write(uint16(addr), ppu.cpu.ram.read(base+uint16(i)))
			addr++
		}

	case 0x4000, 0x4001, 0x4002, 0x4003, 0x4004,
		0x4005, 0x4006, 0x4007, 0x4008, 0x4009, 0x400A,
		0x400B, 0x400C, 0x400D, 0x400E, 0x400F, 0x4010,
		0x4011, 0x4012, 0x4013, 0x4015:
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

func (ppu *Ppu) startVblank() {
	ppu.fVerticalBlank = true
	ppu.cpu.Interrupt(NMI)
}

func (ppu *Ppu) PatternTableAddress(x, y int, baseAddress uint16, tileNumber uint8) uint16 {
	return baseAddress +
		uint16((int(tileNumber)*PATTERN_BYTES_PER_TILE)+
			(y%VERTICAL_PIXES_PER_TILE))
}

func (ppu *Ppu) patternColorIndex(x, y int, baseAddress uint16, tileNumber uint8) uint8 {
	patternAddress := ppu.PatternTableAddress(x, y, baseAddress, tileNumber)
	pattern1 := ppu.vram.read(patternAddress)
	pattern2 := ppu.vram.read(patternAddress + 8)

	bitOffset := uint8((x % 8))
	pattern1 = (pattern1 << bitOffset) >> 7
	pattern2 = ((pattern2 << bitOffset) >> 7)
	return (pattern2 << 1) + pattern1
}

func (ppu *Ppu) backgroundPatternColor() uint8 {
	tileNumber := ppu.nameTablePattern()
	return ppu.patternColorIndex(int(ppu.x), int(ppu.y), ppu.backgroundPatterTableAddress, tileNumber)
}

func (ppu *Ppu) backgroundColorIndex(x, y int) uint8 {

	backgroundColourIndex := uint8(0)

	if ppu.displayBackground && (ppu.showclipBackground || x > 7) {
		patternIndex := ppu.backgroundPatternColor()
		attributeIndex := ppu.nameTableAttribute()
		if patternIndex != 0 {
			backgroundColourIndex = patternIndex + (attributeIndex << 2)
		}
	}

	return backgroundColourIndex
}

func (ppu *Ppu) spriteColorIndex(x, y int, backgroundColorIndex uint8) uint8 {
	index := uint8(0)
	sprite := ppu.sprites[x]

	if !sprite.visible ||
		(sprite.behindBackground && backgroundColorIndex%4 != 0) ||
		(!ppu.showclipSprite && x <= 7) {
		index = ppu.vram.read(0x3F00 + uint16(backgroundColorIndex))
	} else {
		index = ppu.vram.read(0x3F10 + uint16(sprite.paletteIndex))
	}

	if backgroundColorIndex%4 != 0 &&
		sprite.paletteIndex%4 != 0 &&
		sprite.visible &&
		ppu.displayBackground &&
		ppu.displaySprite &&
		(ppu.showclipBackground || x > 7) &&
		(ppu.showclipSprite || x > 7) &&
		!ppu.fSpritZeroHit &&
		x != 255 &&
		y < 239 {
		ppu.fSpritZeroHit = true
	}

	return index
}

func (ppu *Ppu) renderPixel() {
	x := ppu.currentScanlineCycle
	y := ppu.scanline

	colorIndex := ppu.spriteColorIndex(x, y, ppu.backgroundColorIndex(x, y))
	color := colorPalette[colorIndex]
	ppu.gui.DrawPixel(x, y, color[0], color[1], color[2])
}

func (ppu *Ppu) oamGetY(index uint8) byte {
	return ppu.oamRam.read(uint16(4 * index))
}

func (ppu *Ppu) oamGetTile(index uint8) uint8 {
	data := uint8(ppu.oamRam.read(uint16(4*index) + 1))
	if ppu.spriteSize == 8 {
		return data
	}
	// fixme
	return (data >> 1) << 1
}

func (ppu *Ppu) oamGetAttribute(index uint8) byte {
	return ppu.oamRam.read(uint16(4*index) + 2)
}

func (ppu *Ppu) oamGetX(index uint8) byte {
	return ppu.oamRam.read(uint16(4*index) + 3)
}

func (ppu *Ppu) calculateSprites(screenY uint8) {
	ppu.sprites = make([]Sprite, SCREEN_WIDTH)

	if ppu.displaySprite {
		for i := uint8(0); i < 64; i++ {
			if ppu.spriteSize == 8 {
				spriteTopY := uint8(ppu.oamGetY(i)) - 1
				inRange := false
				yOffset := uint8(0)

				if screenY >= spriteTopY {
					yOffset = screenY - spriteTopY
					if yOffset < 8 {
						inRange = true
					}
				}

				if inRange {
					spriteLetfX := ppu.oamGetX(i)
					patternTileNumber := ppu.oamGetTile(i)

					if yOffset >= 8 {
						yOffset -= 8
						patternTileNumber++
					}

					attributeByte := ppu.oamGetAttribute(i)
					flippedHorizontal := (attributeByte>>6)&1 == 1
					flippedVertical := (attributeByte>>7)&1 == 1
					attributeColorIndex := (attributeByte << 6) >> 4

					if flippedVertical {
						yOffset = 8 - yOffset - 1
					}

					for j := uint8(0); j < 8 && (uint16(j)+uint16(spriteLetfX) < SCREEN_WIDTH); j++ {
						x := j
						if flippedHorizontal {
							x = 8 - x - 1
						}

						// todo sprite priority
						if !ppu.sprites[spriteLetfX+j].visible {
							patternColorIndex := ppu.patternColorIndex(int(x), int(yOffset), ppu.spritePatternTableAddress, patternTileNumber)
							if patternColorIndex != 0 {
								sprite := &ppu.sprites[spriteLetfX+j]
								sprite.visible = true
								sprite.paletteIndex = patternColorIndex + attributeColorIndex
								sprite.behindBackground = attributeByte>>5&1 == 1
							}
						}

					}
				}
			} else {
				fatal("8x16 sprite")
			}
		}
	}
}

// scrolling
func (ppu *Ppu) initScroll() {
	ppu.y = uint16(ppu.scrollY)
	ppu.incrementScrollY()
}

func (ppu *Ppu) incrementScrollY() {
	ppu.scrollBase = ppu.basenameTableAddress

	ppu.x = uint16(ppu.scrollX) - 1
	ppu.incrementScrollX()

	ppu.y++

	if ppu.y == TILES_PER_COLUMN*8 {
		switch ppu.scrollBase {
		case 0x2000, 0x2400:
			ppu.scrollBase += 0x0800
		case 0x2800, 0x2C00:
			ppu.scrollBase -= 0x0800
		}
		ppu.y = 0
	}

}

func (ppu *Ppu) incrementScrollX() {
	ppu.x++

	if ppu.x == TILES_PER_ROW*8 {
		switch ppu.scrollBase {
		case 0x2000, 0x2800:
			ppu.scrollBase += 0x0400
		case 0x2400, 0x2C00:
			ppu.scrollBase -= 0x0400
		}

		ppu.x = 0
	}
}

func (ppu *Ppu) tileNumber() uint16 {
	return (ppu.y/8)*TILES_PER_ROW + (ppu.x / 8)
}

func (ppu *Ppu) nameTablePattern() uint8 {
	return ppu.vram.read(ppu.scrollBase + ppu.tileNumber())
}

var attributeTableLookup = [][]uint8{
	{6, 6, 4, 4},
	{6, 6, 4, 4},
	{2, 2, 0, 0},
	{2, 2, 0, 0}}

func (ppu *Ppu) nameTableAttribute() uint8 {
	tileNumber := ppu.tileNumber()
	row := tileNumber / TILES_PER_ROW
	col := tileNumber % TILES_PER_ROW

	attributeByte := ppu.vram.read(ppu.scrollBase + 960 + ((row/4)*8 + col/4))

	return (attributeByte << attributeTableLookup[row%4][col%4]) >> 6
}

func (ppu *Ppu) step() {
	ppu.currentScanlineCycle++

	if ppu.currentScanlineCycle == 341 {
		ppu.scanline++
		if ppu.scanline == 261 {
			ppu.scanline = -1
			ppu.gui.RefreshScreen()
		}
		ppu.currentScanlineCycle = 0

		if ppu.scanline == 241 {
			ppu.startVblank()
		}

		if ppu.scanline == -1 {
			ppu.fVerticalBlank = false
			ppu.fSpriteOverflow = false
			ppu.fSpritZeroHit = false
			ppu.initScroll()
		}

		if ppu.scanline >= -1 && ppu.scanline < SCREEN_HEIGHT {
			ppu.calculateSprites(uint8(ppu.scanline + 1))
			ppu.incrementScrollY()
		}
	}

	if ppu.scanline >= 0 && ppu.scanline < SCREEN_HEIGHT {
		if ppu.currentScanlineCycle >= 0 && ppu.currentScanlineCycle < SCREEN_WIDTH {
			if ppu.currentScanlineCycle > 0 {
				ppu.incrementScrollX()
			}

			ppu.renderPixel()
		}
	}
}

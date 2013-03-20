package main

import (
	"image/png"
	. "launchpad.net/gocheck"
	"os"
	"path/filepath"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func TestPpu(t *testing.T) { TestingT(t) }

type PpuSuite struct{}

var _ = Suite(&PpuSuite{})

func (s *PpuSuite) TestVblank(c *C) {
	ppu := newPpu(&Rom{})
	c.Check(ppu.getStatus(), Equals, uint8(0))
	c.Check(ppu.read(0x2002), Equals, uint8(0))
	ppu.fVerticalBlank = true
	c.Check(ppu.read(0x2002), Equals, uint8(0x80))
	c.Check(ppu.read(0x2002), Equals, uint8(0x00))
	c.Check(ppu.fVerticalBlank, Equals, false)
	c.Check(ppu.getStatus(), Equals, uint8(0x00))
}

func (s *PpuSuite) TestRendering(c *C) {
	filepath.Walk(filepath.Join(fileDir(), "ppu_tests", "bin"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			runRom(path, c)
		}
		return nil
	})
}

func screenshotPath(romPath string) string {
	return filepath.Join(filepath.Dir(filepath.Dir(romPath)), "screenshots", filepath.Base(romPath)+".png")
}

func runRom(romPath string, c *C) {
	rom := LoadRom(romPath)
	ppu := newPpu(rom)
	cpu := newCpu(rom, ppu)
	for i := 0; i < 1000000; i++ {
		run(cpu)
	}
	screenshot := ppu.gui.takeScreenShot()

	path := screenshotPath(romPath)
	previous, err := os.Open(path)
	if err != nil {
		file, _ := os.Create(path)
		png.Encode(file, screenshot)
	} else {
		previousScreenshot, _ := png.Decode(previous)
		for y := 0; y < SCREEN_HEIGHT; y++ {
			for x := 0; x < SCREEN_WIDTH; x++ {
				r, g, b, a := previousScreenshot.At(x, y).RGBA()
				r1, g1, b1, a1 := screenshot.At(x, y).RGBA()
				c.Check(r, Equals, r1)
				c.Check(g, Equals, g1)
				c.Check(b, Equals, b1)
				c.Check(a, Equals, a1)
			}
		}

	}

}

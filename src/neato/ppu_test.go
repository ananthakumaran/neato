package main

import (
	. "launchpad.net/gocheck"
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

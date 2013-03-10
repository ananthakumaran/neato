package main

import (
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func TestMemory(t *testing.T) { TestingT(t) }

type MemorySuite struct{}

var _ = Suite(&MemorySuite{})

func (s *MemorySuite) TestMemoryMirror(c *C) {
	memory := newMemory(0xFFFF)
	memory.mirror(0x0000, 0x07FF, 0x0800, 0x1FFF)
	memory.mirror(0x2000, 0x2007, 0x2008, 0x3FFF)

	memory.write(0x0973, 0x55)
	c.Check(memory.read(0x0973), Equals, uint8(0x55))
	c.Check(memory.read(0x0173), Equals, uint8(0x55))
	c.Check(memory.read(0x1173), Equals, uint8(0x55))
	c.Check(memory.read(0x1973), Equals, uint8(0x55))
	memory.write(0x1173, 0x66)
	c.Check(memory.read(0x0973), Equals, uint8(0x66))
	c.Check(memory.read(0x0173), Equals, uint8(0x66))
	c.Check(memory.read(0x1173), Equals, uint8(0x66))
	c.Check(memory.read(0x1973), Equals, uint8(0x66))

}

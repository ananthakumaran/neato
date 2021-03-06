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
	memory := NewMemory(0xFFFF)
	memory.Mirror(0x0000, 0x07FF, 0x0800, 0x1FFF)
	memory.Mirror(0x2000, 0x2007, 0x2008, 0x3FFF)
	memory.Mirror(0x4000, 0x4FFF, 0xF000, 0xFFFF)

	memory.Write(0x0973, 0x55)
	c.Check(memory.Read(0x0973), Equals, uint8(0x55))
	c.Check(memory.Read(0x0173), Equals, uint8(0x55))
	c.Check(memory.Read(0x1173), Equals, uint8(0x55))
	c.Check(memory.Read(0x1973), Equals, uint8(0x55))
	memory.Write(0x1173, 0x66)
	c.Check(memory.Read(0x0973), Equals, uint8(0x66))
	c.Check(memory.Read(0x0173), Equals, uint8(0x66))
	c.Check(memory.Read(0x1173), Equals, uint8(0x66))
	c.Check(memory.Read(0x1973), Equals, uint8(0x66))
}

func (s *MemorySuite) TestCallbacks(c *C) {
	memory := NewMemory(0xFFFF)
	memory.ReadCallback(0x2000, 0x2007, func(address uint16) byte { return 0x42 })

	lastVal := byte(0)
	memory.WriteCallback(0x2000, 0x2007, func(address uint16, val byte) {
		lastVal = val
	})

	c.Check(memory.Read(0x1FFF), Equals, uint8(0x0))
	c.Check(memory.Read(0x2000), Equals, uint8(0x42))
	memory.Write(0x2000, uint8(0x52))
	c.Check(lastVal, Equals, uint8(0x52))
	c.Check(memory.Read(0x2000), Equals, uint8(0x42))
	c.Check(memory.Read(0x2007), Equals, uint8(0x42))
	c.Check(memory.Read(0x2008), Equals, uint8(0x0))
}

func DummyRead(address uint16) byte {
	return 0
}

func DummyWrite(address uint16, value byte) {

}

func BenchmarkReadWrite(b *testing.B) {
	memory := NewMemory(0xFFFF)
	memory.ReadCallback(0x2000, 0x2007, DummyRead)
	memory.WriteCallback(0x2000, 0x2007, DummyWrite)
	memory.ReadCallback(0x4000, 0x401F, DummyRead)
	memory.WriteCallback(0x4000, 0x401F, DummyWrite)
	memory.Mirror(0x0000, 0x07FF, 0x0800, 0x1FFF)
	memory.Mirror(0x2000, 0x2007, 0x2008, 0x3FFF)

	for i := 0; i < b.N; i++ {
		for j := uint16(0); j < 0xFFFF; j += 0xFF {
			memory.Read(j)
			memory.Write(j, 0)
		}
	}
}

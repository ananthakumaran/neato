package main

import (
	"fmt"
	"neato/petascii"
)

const (
	IRQ   = 0
	NMI   = 1
	RESET = 2
)

var Cycles = [0x100]int{
	7, 6, -1, 8, 3, 3, 5, 5, 3, 2,
	2, 2, 4, 4, 6, 6, 2, 5, -1, 8,
	4, 4, 6, 6, 2, 4, 2, 7, 4, 4,
	7, 7, 6, 6, -1, 8, 3, 3, 5, 5,
	4, 2, 2, 2, 4, 4, 6, 6, 2, 5,
	-1, 8, 4, 4, 6, 6, 2, 4, 2, 7,
	4, 4, 7, 7, 6, 6, -1, 8, 3, 3,
	5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	2, 5, -1, 8, 4, 4, 6, 6, 2, 4,
	2, 7, 4, 4, 7, 7, 6, 6, -1, 8,
	3, 3, 5, 5, 4, 2, 2, 2, 5, 4,
	6, 6, 2, 5, -1, 8, 4, 4, 6, 6,
	2, 4, 2, 7, 4, 4, 7, 7, 2, 6,
	2, 6, 3, 3, 3, 3, 2, 2, 2, 2,
	4, 4, 4, 4, 2, 6, -1, 6, 4, 4,
	4, 4, 2, 5, 2, 5, 5, 5, 5, 5,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2,
	2, 2, 4, 4, 4, 4, 2, 5, -1, 5,
	4, 4, 4, 4, 2, 4, 2, 4, 4, 4,
	4, 4, 2, 6, 2, 8, 3, 3, 5, 5,
	2, 2, 2, 2, 4, 4, 6, 6, 2, 5,
	-1, 8, 4, 4, 6, 6, 2, 4, 2, 7,
	4, 4, 7, 7, 2, 6, 2, 8, 3, 3,
	5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, -1, 8, 4, 4, 6, 6, 2, 4,
	2, 7, 4, 4, 7, 7}

var Bytes = [0x100]int{
	1, 2, -1, 2, 2, 2, 2, 2, 1, 2,
	1, 2, 3, 3, 3, 3, 2, 2, -1, 2,
	2, 2, 2, 2, 1, 3, 1, 3, 3, 3,
	3, 3, -13, 2, -1, 2, 2, 2, 2, 2,
	1, 2, 1, 2, 3, 3, 3, 3, 2, 2,
	-1, 2, 2, 2, 2, 2, 1, 3, 1, 3,
	3, 3, 3, 3, -11, 2, -1, 2, 2, 2,
	2, 2, 1, 2, 1, 2, -13, 3, 3, 3,
	2, 2, -1, 2, 2, 2, 2, 2, 1, 3,
	1, 3, 3, 3, 3, 3, -11, 2, -1, 2,
	2, 2, 2, 2, 1, 2, 1, 2, -13, 3,
	3, 3, 2, 2, -1, 2, 2, 2, 2, 2,
	1, 3, 1, 3, 3, 3, 3, 3, 2, 2,
	2, 2, 2, 2, 2, 2, 1, 2, 1, 2,
	3, 3, 3, 3, 2, 2, -1, 2, 2, 2,
	2, 2, 1, 3, 1, -2, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2, 1, 2,
	1, 2, 3, 3, 3, 3, 2, 2, -1, 2,
	2, 2, 2, 2, 1, 3, 1, 3, 3, 3,
	3, 3, 2, 2, 2, 2, 2, 2, 2, 2,
	1, 2, 1, 2, 3, 3, 3, 3, 2, 2,
	-1, 2, 2, 2, 2, 2, 1, 3, 1, 3,
	3, 3, 3, 3, 2, 2, 2, 2, 2, 2,
	2, 2, 1, 2, 1, 2, 3, 3, 3, 3,
	2, 2, -1, 2, 2, 2, 2, 2, 1, 3,
	1, 3, 3, 3, 3, 3}

var Opcodes = [0x100]string{
	"BRK", "ORA", "KIL", "SLO", "NOP", "ORA", "ASL", "SLO", "PHP", "ORA",
	"ASL", "ANC", "NOP", "ORA", "ASL", "SLO", "BPL", "ORA", "KIL", "SLO",
	"NOP", "ORA", "ASL", "SLO", "CLC", "ORA", "NOP", "SLO", "NOP", "ORA",
	"ASL", "SLO", "JSR", "AND", "KIL", "RLA", "BIT", "AND", "ROL", "RLA",
	"PLP", "AND", "ROL", "ANC", "BIT", "AND", "ROL", "RLA", "BMI", "AND",
	"KIL", "RLA", "NOP", "AND", "ROL", "RLA", "SEC", "AND", "NOP", "RLA",
	"NOP", "AND", "ROL", "RLA", "RTI", "EOR", "KIL", "SRE", "NOP", "EOR",
	"LSR", "SRE", "PHA", "EOR", "LSR", "ALR", "JMP", "EOR", "LSR", "SRE",
	"BVC", "EOR", "KIL", "SRE", "NOP", "EOR", "LSR", "SRE", "CLI", "EOR",
	"NOP", "SRE", "NOP", "EOR", "LSR", "SRE", "RTS", "ADC", "KIL", "RRA",
	"NOP", "ADC", "ROR", "RRA", "PLA", "ADC", "ROR", "ARR", "JMP", "ADC",
	"ROR", "RRA", "BVS", "ADC", "KIL", "RRA", "NOP", "ADC", "ROR", "RRA",
	"SEI", "ADC", "NOP", "RRA", "NOP", "ADC", "ROR", "RRA", "NOP", "STA",
	"NOP", "SAX", "STY", "STA", "STX", "SAX", "DEY", "NOP", "TXA", "XAA",
	"STY", "STA", "STX", "SAX", "BCC", "STA", "KIL", "AHX", "STY", "STA",
	"STX", "SAX", "TYA", "STA", "TXS", "TAS", "SHY", "STA", "SHX", "AHX",
	"LDY", "LDA", "LDX", "LAX", "LDY", "LDA", "LDX", "LAX", "TAY", "LDA",
	"TAX", "LAX", "LDY", "LDA", "LDX", "LAX", "BCS", "LDA", "KIL", "LAX",
	"LDY", "LDA", "LDX", "LAX", "CLV", "LDA", "TSX", "LAS", "LDY", "LDA",
	"LDX", "LAX", "CPY", "CMP", "NOP", "DCP", "CPY", "CMP", "DEC", "DCP",
	"INY", "CMP", "DEX", "AXS", "CPY", "CMP", "DEC", "DCP", "BNE", "CMP",
	"KIL", "DCP", "NOP", "CMP", "DEC", "DCP", "CLD", "CMP", "NOP", "DCP",
	"NOP", "CMP", "DEC", "DCP", "CPX", "SBC", "NOP", "ISC", "CPX", "SBC",
	"INC", "ISC", "INX", "SBC", "NOP", "SBC", "CPX", "SBC", "INC", "ISC",
	"BEQ", "SBC", "KIL", "ISC", "NOP", "SBC", "INC", "ISC", "SED", "SBC",
	"NOP", "ISC", "NOP", "SBC", "INC", "ISC"}

var AddressingMode = [0x100]string{
	"", "izx", "CRASH", "izx", "zp", "zp", "zp", "zp", "", "imm",
	"acc", "imm", "abs", "abs", "abs", "abs", "rel", "izy", "CRASH", "izy",
	"zpx", "zpx", "zpx", "zpx", "", "absy", "", "absy", "absx", "absx",
	"absx", "absx", "abs", "izx", "CRASH", "izx", "zp", "zp", "zp", "zp",
	"", "imm", "acc", "imm", "abs", "abs", "abs", "abs", "rel", "izy",
	"CRASH", "izy", "zpx", "zpx", "zpx", "zpx", "", "absy", "", "absy",
	"absx", "absx", "absx", "absx", "", "izx", "CRASH", "izx", "zp", "zp",
	"zp", "zp", "", "imm", "acc", "imm", "abs", "abs", "abs", "abs",
	"rel", "izy", "CRASH", "izy", "zpx", "zpx", "zpx", "zpx", "", "absy",
	"", "absy", "absx", "absx", "absx", "absx", "", "izx", "CRASH", "izx",
	"zp", "zp", "zp", "zp", "", "imm", "acc", "imm", "ind", "abs",
	"abs", "abs", "rel", "izy", "CRASH", "izy", "zpx", "zpx", "zpx", "zpx",
	"", "absy", "", "absy", "absx", "absx", "absx", "absx", "imm", "izx",
	"imm", "izx", "zp", "zp", "zp", "zp", "", "imm", "", "imm",
	"abs", "abs", "abs", "abs", "rel", "izy", "CRASH", "izy", "zpx", "zpx",
	"zpy", "zpy", "", "absy", "", "", "absx", "absx", "absy", "absy",
	"imm", "izx", "imm", "izx", "zp", "zp", "zp", "zp", "", "imm",
	"", "imm", "abs", "abs", "abs", "abs", "rel", "izy", "CRASH", "izy",
	"zpx", "zpx", "zpy", "zpy", "", "absy", "", "absy", "absx", "absx",
	"absy", "absy", "imm", "izx", "imm", "izx", "zp", "zp", "zp", "zp",
	"", "imm", "", "imm", "abs", "abs", "abs", "abs", "rel", "izy",
	"CRASH", "izy", "zpx", "zpx", "zpx", "zpx", "", "absy", "", "absy",
	"absx", "absx", "absx", "absx", "imm", "izx", "imm", "izx", "zp", "zp",
	"zp", "zp", "", "imm", "", "imm", "abs", "abs", "abs", "abs",
	"rel", "izy", "CRASH", "izy", "zpx", "zpx", "zpx", "zpx", "", "absy",
	"", "absy", "absx", "absx", "absx"}

type Cpu struct {
	rom *Rom
	ppu *Ppu

	ram *Memory

	// registers
	pc    uint16
	ac    uint8
	x     uint8
	y     uint8
	stack uint8

	// status flags
	fCarry            bool
	fZero             bool
	fInterruptDisable bool
	fDecimal          bool
	fBreak            bool
	fNotUsed          bool
	fOverflow         bool
	fNegative         bool

	// output
	printer *petascii.Printer

	// interrupts
	pendingInterruptRequest bool
	interruptType           int

	cycles int

	// behavior
	respectDecimalMode bool
}

func newCpu(rom *Rom, ppu *Ppu) *Cpu {
	cpu := Cpu{}
	cpu.rom = rom
	cpu.ram = newMemory(0xFFFF)
	cpu.ram.copy(0x8000, 0xC000, rom.PrgRoms[0])

	switch rom.PrgRomCount {
	case 1:
		cpu.ram.copy(0xC000, 0x10000, rom.PrgRoms[0])
	case 2:
		cpu.ram.copy(0xC000, 0x10000, rom.PrgRoms[1])
	default:
		fatal("uknown prg rom count")
	}

	// IO registers
	cpu.ram.readCallback(0x2000, 0x2007, func(address uint16) byte { return cpu.ppu.read(address) })
	cpu.ram.writeCallback(0x2000, 0x2007, func(address uint16, val byte) { cpu.ppu.write(address, val) })
	cpu.ram.readCallback(0x4000, 0x401F, func(address uint16) byte { return cpu.ppu.read(address) })
	cpu.ram.writeCallback(0x4000, 0x401F, func(address uint16, val byte) { cpu.ppu.write(address, val) })

	cpu.ram.mirror(0x0000, 0x07FF, 0x0800, 0x1FFF)
	cpu.ram.mirror(0x2000, 0x2007, 0x2008, 0x3FFF)

	cpu.ppu = ppu
	ppu.cpu = &cpu
	cpu.respectDecimalMode = false
	cpu.reset()
	return &cpu
}

func (cpu *Cpu) step() int {

	cpu.handleInterrupt()

	ir := cpu.read(cpu.pc)
	address := uint16(0)
	immediate := false
	accumulator := false
	// info("%02X  %02d ", ir, ir)

	info("%04X  ", cpu.pc)

	bytes := Bytes[ir]
	if bytes <= 0 && bytes >= -10 {
		fmt.Printf(" ir %x %d bytes %d", ir, ir, bytes)
		fatal("invalid bytes")
	}

	if bytes <= -10 {
		bytes = -bytes - 10
	}

	for i := 0; i < 3; i++ {
		if bytes > 0 {
			info("%02X ", cpu.read(cpu.pc+uint16(i)))
			bytes--
		} else {
			info("   ")
		}
	}

	if Opcodes[ir] == "NOP" && ir != 0xEA {
		info("*")
	} else {
		info(" ")
	}

	info("%s ", Opcodes[ir])
	cpu.cycles = Cycles[ir]

	switch AddressingMode[ir] {
	case "abs":
		address = uint16(cpu.read(cpu.pc+2))<<8 | uint16(cpu.read(cpu.pc+1))

		switch Opcodes[ir] {
		case "JMP", "JSR":
			// case "STX", "LDX", "STY", "LDY", "STA", "LDA",
			// 	"BIT", "ORA":
			info("$%04X                       ", address)
		default:
			info("$%04X = %02X                  ", address, cpu.dummyRead(address))

		}

	case "absx":
		tempAddress := uint16(cpu.read(cpu.pc+2))<<8 | uint16(cpu.read(cpu.pc+1))
		address = tempAddress + uint16(cpu.x)
		info("$%04X,X @ %04X = %02X         ", tempAddress, address, cpu.dummyRead(address))
		switch Opcodes[ir] {
		case "ADC", "AND", "CMP", "EOR", "LDA", "LDY", "ORA", "SBC", "NOP":
			if tempAddress>>8 != address>>8 {
				cpu.cycles++
			}
		}

	case "absy":
		tempAddress := uint16(cpu.read(cpu.pc+2))<<8 | uint16(cpu.read(cpu.pc+1))
		address = tempAddress + uint16(cpu.y)
		info("$%04X,Y @ %04X = %02X         ", tempAddress, address, cpu.dummyRead(address))

		if Opcodes[ir] != "STA" && tempAddress>>8 != address>>8 {
			cpu.cycles++
		}

	case "zp":
		address = uint16(cpu.read(cpu.pc + 1))
		info("$%02X = %02X                    ", cpu.read(cpu.pc+1), cpu.dummyRead(address))

	case "zpx":
		base := cpu.read(cpu.pc + 1)
		address = uint16(base + cpu.x)
		info("$%02X,X @ %02X = %02X             ", base, address, cpu.dummyRead(address))

	case "zpy":
		base := cpu.read(cpu.pc + 1)
		address = uint16(base + cpu.y)
		info("$%02X,Y @ %02X = %02X             ", base, address, cpu.dummyRead(address))

	case "izy":
		base := (cpu.read(cpu.pc + 1))
		tempAddress := uint16(cpu.read(uint16(base+1)))<<8 | uint16(cpu.read(uint16(base)))
		address = tempAddress + uint16(cpu.y)

		if (tempAddress >> 8) != (address >> 8) {
			cpu.cycles++
		}

		info("($%02X),Y = %04X @ %04X = %02X  ", base, tempAddress, address, cpu.dummyRead(address))

	case "izx":
		val := cpu.read(cpu.pc + 1)
		base := val + cpu.x
		address = uint16(cpu.read(uint16(base+1)))<<8 | uint16(cpu.read(uint16(base)))
		info("($%02X,X) @ %02X = %04X = %02X    ", val, base, address, cpu.dummyRead(address))

	case "ind":
		base := uint16(cpu.read(cpu.pc+2))<<8 | uint16(cpu.read(cpu.pc+1))
		address = uint16(cpu.read(base+1))<<8 | uint16(cpu.read(base))
		info("($%04X) = %04X              ", base, address)

		// page boundary bug
		if base&0x00FF == 0x00FF {
			address = uint16(cpu.read(base&0xFF00))<<8 | uint16(cpu.read(base))
		}

	case "imm":
		immediate = true
		info("#$%02X                        ", cpu.read(cpu.pc+1))
	case "acc":
		accumulator = true
		info("A                           ")
	case "":
		info("                            ")
	case "rel": // never mind
		address = uint16(int(cpu.pc)+int(int8(cpu.read(cpu.pc+1)))) + uint16(Bytes[ir])
		info("$%04X                       ", address)
	default:
		fatal("unknown addressing mode not implemented %x", AddressingMode[ir])
	}

	info("A:%02X X:%02X Y:%02X P:%02X SP:%02X", cpu.ac, cpu.x, cpu.y, cpu.getStatus(), cpu.stack)

	switch Opcodes[ir] {
	case "CLD":
		cpu.fDecimal = false
	case "CLC":
		cpu.fCarry = false
	case "SEC":
		cpu.fCarry = true
	case "SEI":
		cpu.fInterruptDisable = true
	case "SED":
		cpu.fDecimal = true
	case "CLI":
		cpu.fInterruptDisable = false
	case "CLV":
		cpu.fOverflow = false
	case "LDA":
		cpu.ac = cpu.val(immediate, address)
		cpu.zeroNeg(cpu.ac)
	case "LDX":
		cpu.x = cpu.val(immediate, address)
		cpu.zeroNeg(cpu.x)
	case "LDY":
		cpu.y = cpu.val(immediate, address)
		cpu.zeroNeg(cpu.y)
	case "CMP":
		val := cpu.val(immediate, address)
		cpu.fCarry = cpu.ac >= val
		cpu.zeroNeg(cpu.ac - val)
	case "CPY":
		val := cpu.val(immediate, address)
		cpu.fCarry = cpu.y >= val
		cpu.zeroNeg(cpu.y - val)
	case "CPX":
		val := cpu.val(immediate, address)
		cpu.fCarry = cpu.x >= val
		cpu.zeroNeg(cpu.x - val)
	case "ADC":
		val := cpu.val(immediate, address)
		if cpu.respectDecimalMode && cpu.fDecimal {
			cpu.fZero = (int(val&0xff)+int(cpu.ac&0xff)+int(cpu.carry()))&0xff == 0
			al := int(cpu.ac&0x0F) + int(val&0x0F) + int(cpu.carry())
			if al >= 0x0A {
				al = (al+0x06)&0x0F + 0x10
			}

			a := int(cpu.ac&0xF0) + int(val&0xF0) + al
			sa := int(int8(cpu.ac&0xF0)) + int(int8(val&0xF0)) + int(al)
			cpu.fNegative = a>>7&1 == 1
			cpu.fOverflow = sa < -128 || sa > 127
			if a >= 0xA0 {
				a += 0x60
			}
			cpu.fCarry = a >= 0x100
			cpu.ac = uint8(a & 0xFF)

		} else {
			result := int(val&0xff) + int(cpu.ac&0xff) + int(cpu.carry())
			carry6 := int(val&0x7f) + int(cpu.ac&0x7f) + int(cpu.carry())
			cpu.fCarry = result&0x100 != 0
			cpu.fOverflow = Xor(cpu.fCarry, ((carry6 & 0x80) != 0))
			cpu.ac = uint8(result & 0xff)
			cpu.zeroNeg(cpu.ac)
		}

	case "SBC":
		val := cpu.val(immediate, address)
		a := 0

		if cpu.respectDecimalMode && cpu.fDecimal {
			al := int(cpu.ac&0x0F) - int(val&0x0f) + (int(cpu.carry()) - 1)
			if al < 0 {
				al = ((al - 0x06) & 0x0F) - 0x10
			}
			a = int(cpu.ac&0xF0) - int(val&0xF0) + al
			if a < 0 {
				a -= 0x60
			}
		}

		val = ^val
		result := int(val)&0xff + int(cpu.ac)&0xff + int(cpu.carry())
		carry6 := int(val)&0x7f + int(cpu.ac)&0x7f + int(cpu.carry())
		cpu.fCarry = result&0x100 != 0
		cpu.fOverflow = Xor(cpu.fCarry, ((carry6 & 0x80) != 0))
		result &= 0xff
		cpu.fZero = result == 0
		cpu.fNegative = result&0x80 != 0

		if cpu.respectDecimalMode && cpu.fDecimal {
			cpu.ac = uint8(a & 0xff)
		} else {
			cpu.ac = uint8(result)
		}

	case "ASL":
		if accumulator {
			cpu.fCarry = (cpu.ac>>7)&1 == 1
			cpu.ac = cpu.ac << 1
			cpu.zeroNeg(cpu.ac)
		} else {
			val := cpu.val(immediate, address)
			cpu.fCarry = (val>>7)&1 == 1
			val = val << 1
			cpu.write(address, val)
			cpu.zeroNeg(val)
		}

	case "LSR":
		if accumulator {
			cpu.fCarry = (cpu.ac & 0x01) == 1
			cpu.ac = cpu.ac >> 1
			cpu.zeroNeg(cpu.ac)
		} else {
			val := cpu.val(immediate, address)
			cpu.fCarry = (val & 0x01) == 1
			val = val >> 1
			cpu.write(address, val)
			cpu.zeroNeg(val)
		}

	case "ROL":
		c := cpu.carry()
		if accumulator {
			cpu.fCarry = (cpu.ac >> 7) == 1
			cpu.ac = cpu.ac<<1 | c
			cpu.zeroNeg(cpu.ac)
		} else {
			val := cpu.val(immediate, address)
			cpu.fCarry = (val >> 7) == 1
			val = val<<1 | c
			cpu.write(address, val)
			cpu.zeroNeg(val)
		}
	case "ROR":
		c := cpu.carry() << 7
		if accumulator {
			cpu.fCarry = (cpu.ac & 0x01) == 1
			cpu.ac = cpu.ac>>1 | c
			cpu.zeroNeg(cpu.ac)
		} else {
			val := cpu.val(immediate, address)
			cpu.fCarry = (val & 0x01) == 1
			val = val>>1 | c
			cpu.write(address, val)
			cpu.zeroNeg(val)
		}

	case "BPL":
		if !cpu.fNegative {
			cpu.relativeJmp()
		}
	case "BMI":
		if cpu.fNegative {
			cpu.relativeJmp()
		}
	case "BNE":
		if !cpu.fZero {
			cpu.relativeJmp()
		}
	case "BEQ":
		if cpu.fZero {
			cpu.relativeJmp()
		}
	case "BCC":
		if !cpu.fCarry {
			cpu.relativeJmp()
		}
	case "BCS":
		if cpu.fCarry {
			cpu.relativeJmp()
		}
	case "BVC":
		if !cpu.fOverflow {
			cpu.relativeJmp()
		}
	case "BVS":
		if cpu.fOverflow {
			cpu.relativeJmp()
		}
	case "JSR":
		cpu.push(uint8((cpu.pc + 2) >> 8))
		cpu.push(uint8((cpu.pc + 2) & 0xFF))
		cpu.pc = address
	case "RTS":
		lo := cpu.pull()
		hi := cpu.pull()
		cpu.pc = uint16(hi)<<8 | uint16(lo)
		cpu.pc += 1
	case "JMP":
		cpu.pc = address
	case "STX":
		cpu.write(address, cpu.x)
	case "STY":
		cpu.write(address, cpu.y)
	case "STA":
		cpu.write(address, cpu.ac)
	case "INX":
		cpu.x += 1
		cpu.zeroNeg(cpu.x)
	case "INY":
		cpu.y += 1
		cpu.zeroNeg(cpu.y)
	case "DEX":
		cpu.x -= 1
		cpu.zeroNeg(cpu.x)
	case "DEY":
		cpu.y -= 1
		cpu.zeroNeg(cpu.y)
	case "INC":
		temp := cpu.val(immediate, address) + 1
		cpu.zeroNeg(temp)
		cpu.write(address, temp)
	case "DEC":
		temp := cpu.val(immediate, address) - 1
		cpu.zeroNeg(temp)
		cpu.write(address, temp)
	case "TXS":
		cpu.stack = cpu.x
	case "TSX":
		cpu.x = cpu.stack
		cpu.zeroNeg(cpu.x)
	case "TXA":
		cpu.ac = cpu.x
		cpu.zeroNeg(cpu.ac)
	case "TAX":
		cpu.x = cpu.ac
		cpu.zeroNeg(cpu.x)
	case "TAY":
		cpu.y = cpu.ac
		cpu.zeroNeg(cpu.y)
	case "TYA":
		cpu.ac = cpu.y
		cpu.zeroNeg(cpu.ac)
	case "EOR":
		cpu.ac = cpu.ac ^ cpu.val(immediate, address)
		cpu.zeroNeg(cpu.ac)
	case "ORA":
		cpu.ac = cpu.ac | cpu.val(immediate, address)
		cpu.zeroNeg(cpu.ac)
	case "AND":
		cpu.ac = cpu.ac & cpu.val(immediate, address)
		cpu.zeroNeg(cpu.ac)
	case "BIT":
		val := cpu.val(immediate, address)
		cpu.fZero = cpu.ac&val == 0
		cpu.fOverflow = (val>>6)&1 == 1
		cpu.fNegative = (val>>7)&1 == 1
	case "PLA":
		cpu.ac = cpu.pull()
		cpu.zeroNeg(cpu.ac)
	case "PLP":
		cpu.status(cpu.pull())
		// cpu.fBreak = false
	case "PHP":
		// check this behaviour
		cpu.fBreak = true
		cpu.push(cpu.getStatus())
		cpu.fBreak = false
	case "PHA":
		cpu.push(cpu.ac)
	case "NOP":
	case "RTI":
		cpu.status(cpu.pull())
		lo := cpu.pull()
		hi := cpu.pull()
		cpu.pc = uint16(hi)<<8 | uint16(lo)
		info("\npoping address %04X\n", cpu.pc)
	case "BRK":
		cpu.push(uint8((cpu.pc + 1) >> 8))
		cpu.push(uint8((cpu.pc + 1) & 0xFF))
		cpu.push(cpu.getStatus())
		cpu.pc = uint16(cpu.read(0xFFFF))<<8 | uint16(cpu.read(0xFFFE))
		cpu.fBreak = true
		//fatal("\nbreak\n")
	default:
		fmt.Printf("ir 0x%X mode %s cycles %d opcode %s byte %d addr 0x%X  ", ir, AddressingMode[ir], Cycles[ir], Opcodes[ir], Bytes[ir], cpu.pc)
		fatal("not implemented")
	}

	if Bytes[ir] > 0 {
		cpu.pc += uint16(Bytes[ir])
	}

	//fmt.Printf("opcode %s ir 0x%X mode %s cycles %d  byte %d addr 0x%X \n ", Opcodes[ir], ir, AddressingMode[ir], Cycles[ir], Bytes[ir], address)
	//cpu.inspect()
	if cpu.cycles <= 0 {
		fatal("invalid cycle %d", cpu.cycles)
	}

	return cpu.cycles
}

func (cpu *Cpu) read(address uint16) byte {
	return cpu.ram.read(address)

}

func (cpu *Cpu) dummyRead(address uint16) byte {
	return 0

}

func (cpu *Cpu) write(address uint16, val uint8) {
	cpu.ram.write(address, val)
}

func (cpu *Cpu) reset() {
	// OxFFFC & 0xFFFD contains the intital PC register
	cpu.pc = (uint16(cpu.ram.read(0xFFFD)) << 8) | uint16(cpu.ram.read(0xFFFC))
	cpu.stack = 0xFD
	cpu.fInterruptDisable = true
	//cpu.pc = 0xC000
}

func (cpu *Cpu) zeroNeg(val uint8) {
	cpu.fZero = val == 0
	cpu.fNegative = val>>7 == 1
}

func (cpu *Cpu) push(val uint8) {
	cpu.write(0x0100|uint16(cpu.stack), val)
	cpu.stack -= 1
}

func (cpu *Cpu) pull() uint8 {
	cpu.stack += 1
	return cpu.read(0x0100 | uint16(cpu.stack))
}

func (cpu *Cpu) val(immediate bool, address uint16) uint8 {
	if address == 0 && !immediate {
		//fatal("invalid address")
	}

	var val byte
	if immediate {
		val = cpu.read(cpu.pc + 1)
	} else {
		val = cpu.read(address)
	}

	return val
}

func (cpu *Cpu) relativeJmp() {
	cpu.cycles++
	address := uint16(int(cpu.pc) + int(int8(cpu.read(cpu.pc+1))))
	if (address >> 8) != (cpu.pc >> 8) {
		// todo check
		//cpu.cycles++
	}
	cpu.pc = address
}

func (cpu *Cpu) carry() uint8 {
	if cpu.fCarry {
		return 1
	}
	return 0
}

// status byte
// 7 6 5 4 3 2 1 0
// | | | | | | | -- carry
// | | | | | | -- zero
// | | | | | -- interrupt disable
// | | | | -- decimal
// | | | -- break
// | | -- unused (always 1)
// | -- overflow
// -- negative
func (cpu *Cpu) status(val uint8) {
	cpu.fCarry = val&1 == 1
	cpu.fZero = (val>>1)&1 == 1
	cpu.fInterruptDisable = (val>>2)&1 == 1
	cpu.fDecimal = (val>>3)&1 == 1
	cpu.fBreak = (val>>4)&1 == 1
	cpu.fNotUsed = (val>>5)&1 == 1
	cpu.fOverflow = (val>>6)&1 == 1
	cpu.fNegative = (val>>7)&1 == 1
}

func (cpu *Cpu) getStatus() uint8 {
	status := uint8(0)
	if cpu.fCarry {
		status |= 0x01
	}

	if cpu.fZero {
		status |= 0x02
	}

	if cpu.fInterruptDisable {
		status |= 0x04
	}

	if cpu.fDecimal {
		status |= 0x08
	}

	if cpu.fBreak {
		status |= 0x10
	}

	status |= 0x20

	if cpu.fOverflow {
		status |= 0x40
	}

	if cpu.fNegative {
		status |= 0x80
	}

	return status
}

func (cpu *Cpu) inspect() {
	debug("pc: %X, ac: %X , x: %X, y: %X, stack: %X status: %b\n", cpu.pc, cpu.ac, cpu.x, cpu.y, cpu.stack, cpu.getStatus())
}

func Xor(a, b bool) bool {
	return (a || b) && !(a && b)
}

func (cpu *Cpu) handleInterrupt() {
	if cpu.pendingInterruptRequest {
		switch cpu.interruptType {
		case IRQ:
			if !cpu.fInterruptDisable {
				cpu.push(uint8((cpu.pc) >> 8))
				cpu.push(uint8((cpu.pc) & 0xFF))
				cpu.push(cpu.getStatus())
				cpu.pc = uint16(cpu.read(0xFFFF))<<8 | uint16(cpu.read(0xFFFE))

			}
		case NMI:
			if cpu.ppu.nmiOnVBlank {
				info("vblank on nmi")
				info("\npushing address %04X\n", cpu.pc+1)
				cpu.push(uint8((cpu.pc) >> 8))
				cpu.push(uint8((cpu.pc) & 0xFF))
				cpu.push(cpu.getStatus())
				cpu.pc = uint16(cpu.read(0xFFFB))<<8 | uint16(cpu.read(0xFFFA))
			}
		case RESET:
			cpu.pc = uint16(cpu.read(0xFFFD))<<8 | uint16(cpu.read(0xFFFC))
		}

		cpu.fInterruptDisable = true
		cpu.pendingInterruptRequest = false
	}
}

func (cpu *Cpu) Interrupt(interruptType int) {
	cpu.pendingInterruptRequest = true
	cpu.interruptType = interruptType
}

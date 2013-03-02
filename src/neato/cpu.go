package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Cycles = [0x100]int{
	0, 6, -1, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
	3, 5, -1, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, -1, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
	2, 5, -1, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, -1, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	3, 5, -1, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, -1, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
	2, 5, -1, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	3, 6, -1, 6, 4, 4, 4, 4, 2, 5, 2, 5, 5, 5, 5, 5,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 5, -1, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	3, 5, -1, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, -1, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7}

var Bytes = [0x100]int{
	0, 2, -1, 2, 2, 2, 2, 2, 1, 2, 1, 2, 3, 3, 3, 3,
	2, 2, -1, 2, 2, 2, 2, 2, 1, 3, 1, 3, 3, 3, 3, 3,
	-2, 2, -1, 2, 2, 2, 2, 2, 1, 2, 1, 2, 3, 3, 3, 3,
	2, 2, -1, 2, 2, 2, 2, 2, 1, 3, 1, 3, 3, 3, 3, 3,
	-2, 2, -1, 2, 2, 2, 2, 2, 1, 2, 1, 2, -2, 3, 3, 3,
	2, 2, -1, 2, 2, 2, 2, 2, 1, 3, 1, 3, 3, 3, 3, 3,
	-2, 2, -1, 2, 2, 2, 2, 2, 1, 2, 1, 2, -2, 3, 3, 3,
	2, 2, -1, 2, 2, 2, 2, 2, 1, 3, 1, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2, 1, 2, 1, 2, 3, 3, 3, 3,
	2, 2, -1, 2, 2, 2, 2, 2, 1, 3, 1, -2, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2, 1, 2, 1, 2, 3, 3, 3, 3,
	2, 2, -1, 2, 2, 2, 2, 2, 1, 3, 1, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2, 1, 2, 1, 2, 3, 3, 3, 3,
	2, 2, -1, 2, 2, 2, 2, 2, 1, 3, 1, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2, 1, 2, 1, 2, 3, 3, 3, 3,
	2, 2, -1, 2, 2, 2, 2, 2, 1, 3, 1, 3, 3, 3, 3, 3}

var Opcodes = [0x100]string{
	"BRK", "ORA", "KIL", "SLO", "NOP", "ORA", "ASL", "SLO",
	"PHP", "ORA", "ASL", "ANC", "NOP", "ORA", "ASL", "SLO",
	"BPL", "ORA", "KIL", "SLO", "NOP", "ORA", "ASL", "SLO",
	"CLC", "ORA", "NOP", "SLO", "NOP", "ORA", "ASL", "SLO",
	"JSR", "AND", "KIL", "RLA", "BIT", "AND", "ROL", "RLA",
	"PLP", "AND", "ROL", "ANC", "BIT", "AND", "ROL", "RLA",
	"BMI", "AND", "KIL", "RLA", "NOP", "AND", "ROL", "RLA",
	"SEC", "AND", "NOP", "RLA", "NOP", "AND", "ROL", "RLA",
	"RTI", "EOR", "KIL", "SRE", "NOP", "EOR", "LSR", "SRE",
	"PHA", "EOR", "LSR", "ALR", "JMP", "EOR", "LSR", "SRE",
	"BVC", "EOR", "KIL", "SRE", "NOP", "EOR", "LSR", "SRE",
	"CLI", "EOR", "NOP", "SRE", "NOP", "EOR", "LSR", "SRE",
	"RTS", "ADC", "KIL", "RRA", "NOP", "ADC", "ROR", "RRA",
	"PLA", "ADC", "ROR", "ARR", "JMP", "ADC", "ROR", "RRA",
	"BVS", "ADC", "KIL", "RRA", "NOP", "ADC", "ROR", "RRA",
	"SEI", "ADC", "NOP", "RRA", "NOP", "ADC", "ROR", "RRA",
	"NOP", "STA", "NOP", "SAX", "STY", "STA", "STX", "SAX",
	"DEY", "NOP", "TXA", "XAA", "STY", "STA", "STX", "SAX",
	"BCC", "STA", "KIL", "AHX", "STY", "STA", "STX", "SAX",
	"TYA", "STA", "TXS", "TAS", "SHY", "STA", "SHX", "AHX",
	"LDY", "LDA", "LDX", "LAX", "LDY", "LDA", "LDX", "LAX",
	"TAY", "LDA", "TAX", "LAX", "LDY", "LDA", "LDX", "LAX",
	"BCS", "LDA", "KIL", "LAX", "LDY", "LDA", "LDX", "LAX",
	"CLV", "LDA", "TSX", "LAS", "LDY", "LDA", "LDX", "LAX",
	"CPY", "CMP", "NOP", "DCP", "CPY", "CMP", "DEC", "DCP",
	"INY", "CMP", "DEX", "AXS", "CPY", "CMP", "DEC", "DCP",
	"BNE", "CMP", "KIL", "DCP", "NOP", "CMP", "DEC", "DCP",
	"CLD", "CMP", "NOP", "DCP", "NOP", "CMP", "DEC", "DCP",
	"CPX", "SBC", "NOP", "ISC", "CPX", "SBC", "INC", "ISC",
	"INX", "SBC", "NOP", "SBC", "CPX", "SBC", "INC", "ISC",
	"BEQ", "SBC", "KIL", "ISC", "NOP", "SBC", "INC", "ISC",
	"SED", "SBC", "NOP", "ISC", "NOP", "SBC", "INC", "ISC"}

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
	rom Rom
	ppu Ppu

	lastPc uint16
	ram    []byte

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
}

func newCpu(rom Rom, ppu Ppu) Cpu {
	cpu := Cpu{}
	cpu.rom = rom
	cpu.ram = make([]byte, 0x10000)
	copy(cpu.ram[0x8000:0xC000], rom.PrgRoms[0])
	copy(cpu.ram[0xC000:0x10000], rom.PrgRoms[0])

	cpu.ppu = ppu
	cpu.reset()
	return cpu
}

func testCpu(filename string) Cpu {
	cpu := Cpu{}
	cpu.ram = make([]byte, 0x10000)
	cpu.loadFile(filename)
	return cpu
}

func (cpu *Cpu) loadFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fatal("file not found: ", filename)
	}

	add := make([]byte, 2)
	count, err := file.Read(add)
	if err != nil || count != 2 {
		fatal("invalid file")
	}
	start := uint16(add[1])<<8 + uint16(add[0])

	count, err = file.Read(cpu.ram[start:])
	if err != nil || count == 0 {
		fatal("invalid file")
	}

	cpu.ram[0x0002] = 0x00
	cpu.ram[0xA002] = 0x00
	cpu.ram[0xA003] = 0x80
	cpu.ram[0xFFFE] = 0x48
	cpu.ram[0xFFFF] = 0xFF
	cpu.ram[0x01FE] = 0xFF
	cpu.ram[0x01FF] = 0x7F

	// FF48  48        PHA
	// FF49  8A        TXA
	// FF4A  48        PHA
	// FF4B  98        TYA
	// FF4C  48        PHA
	// FF4D  BA        TSX
	// FF4E  BD 04 01  LDA    $0104,X
	// FF51  29 10     AND    #$10
	// FF53  F0 03     BEQ    $FF58
	// FF55  6C 16 03  JMP    ($0316)
	// FF58  6C 14 03  JMP    ($0314)

	copy(cpu.ram[0xFF48:0xFF5A], []byte{0x48, 0x8A, 0x48, 0x98, 0x48, 0xBA,
		0xBD, 0x04, 0x01, 0x29, 0x10, 0xF0,
		0x03, 0x6C, 0x16, 0x03, 0x6C, 0x14, 0x03})

	cpu.status(0x04)
	cpu.stack = 0xFD
	cpu.pc = 0x0801
}

func (cpu *Cpu) step() {
	cpu.lastPc = cpu.pc

	switch cpu.pc {
	case 0xFFD2:
		cpu.ram[0x030C] = 0
		if cpu.ac == 0x0D {
			fmt.Printf("\n")
		} else if strconv.IsPrint(rune(cpu.ac)) {
			fmt.Printf("%c", cpu.ac)
		} else {
			fmt.Printf("-%X-", cpu.ac)
		}
		lo := cpu.pull()
		hi := cpu.pull()
		cpu.pc = uint16(hi)<<8 | uint16(lo)
		cpu.pc += 1
		return

	case 0xE16F:
		os.Exit(0)
		lo := cpu.read(0xBB)
		hi := cpu.read(0xBC)
		start := uint16(hi)<<8 | uint16(lo)
		length := cpu.read(0xB7)
		filename := string(cpu.ram[start : start+uint16(length)])
		fmt.Println("\n", filename)
		cpu.loadFile("suite/bin/" + strings.ToLower(filename))
		return

	case 0xFFE4:
		cpu.ac = 3
		lo := cpu.pull()
		hi := cpu.pull()
		cpu.pc = uint16(hi)<<8 | uint16(lo)
		cpu.pc += 1
		return

	case 0x8000, 0xA474:
		fatal("exit trap")
	}

	ir := cpu.read(cpu.pc)
	address := uint16(0)
	immediate := false
	accumulator := false

	switch AddressingMode[ir] {
	case "abs":
		address = uint16(cpu.read(cpu.pc+2))<<8 | uint16(cpu.read(cpu.pc+1))
	case "absx":
		address = uint16(cpu.read(cpu.pc+2))<<8 | uint16(cpu.read(cpu.pc+1))
		address += uint16(cpu.x)
	case "absy":
		address = uint16(cpu.read(cpu.pc+2))<<8 | uint16(cpu.read(cpu.pc+1))
		address += uint16(cpu.y)
	case "zp":
		address = uint16(cpu.read(cpu.pc + 1))

	case "zpx":
		address = uint16(cpu.read(cpu.pc+1) + cpu.x)

	case "zpy":
		address = uint16(cpu.read(cpu.pc+1) + cpu.y)

	case "izy":
		base := uint16((cpu.read(cpu.pc + 1)))
		address = uint16(cpu.read(base+1))<<8 | uint16(cpu.read(base))
		address += uint16(cpu.y)

	case "izx":
		val := cpu.read(cpu.pc + 1)
		base := uint16(val + cpu.x)
		address = uint16(cpu.read(base+1))<<8 | uint16(cpu.read(base))

	case "ind":
		base := uint16((cpu.read(cpu.pc + 1)))
		address = uint16(cpu.read(base+1))<<8 | uint16(cpu.read(base))
		address = uint16(cpu.read(address+1))<<8 | uint16(cpu.read(address))
	case "imm":
		immediate = true
	case "acc":
		accumulator = true
	case "", "rel": // never mind
	default:
		fatal("unknown addressing mode not implemented %x", AddressingMode[ir])
	}

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
		if cpu.fDecimal {
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

		if cpu.fDecimal {
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

		if cpu.fDecimal {
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
			cpu.pc = uint16(int(cpu.pc) + int(int8(cpu.read(cpu.pc+1))))
		}
	case "BMI":
		if cpu.fNegative {
			cpu.pc = uint16(int(cpu.pc) + int(int8(cpu.read(cpu.pc+1))))
		}
	case "BNE":
		if !cpu.fZero {
			cpu.pc = uint16(int(cpu.pc) + int(int8(cpu.read(cpu.pc+1))))
		}
	case "BEQ":
		if cpu.fZero {
			cpu.pc = uint16(int(cpu.pc) + int(int8(cpu.read(cpu.pc+1))))
		}
	case "BCC":
		if !cpu.fCarry {
			cpu.pc = uint16(int(cpu.pc) + int(int8(cpu.read(cpu.pc+1))))
		}
	case "BCS":
		if cpu.fCarry {
			cpu.pc = uint16(int(cpu.pc) + int(int8(cpu.read(cpu.pc+1))))
		}
	case "BVC":
		if !cpu.fOverflow {
			cpu.pc = uint16(int(cpu.pc) + int(int8(cpu.read(cpu.pc+1))))
		}
	case "BVS":
		if cpu.fOverflow {
			cpu.pc = uint16(int(cpu.pc) + int(int8(cpu.read(cpu.pc+1))))
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
		// page boundary fix
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
	case "PHP":
		cpu.fBreak = true
		cpu.push(cpu.getStatus())
	case "PHA":
		cpu.push(cpu.ac)
	case "NOP":
	case "RTI":
		cpu.status(cpu.pull())
		lo := cpu.pull()
		hi := cpu.pull()
		cpu.pc = uint16(hi)<<8 | uint16(lo)
	case "BRK":
		cpu.push(uint8((cpu.pc + 1) >> 8))
		cpu.push(uint8((cpu.pc + 1) & 0xFF))
		cpu.push(cpu.getStatus())
		cpu.pc = uint16(cpu.read(0xFFFF))<<8 | uint16(cpu.read(0xFFFE))
		cpu.fBreak = true
		fatal("\nbreak\n")
	default:
		fmt.Printf("ir 0x%X mode %s cycles %d opcode %s byte %d addr 0x%X  ", ir, AddressingMode[ir], Cycles[ir], Opcodes[ir], Bytes[ir], cpu.pc)
		fatal("not implemented")
	}

	if Bytes[ir] > 0 {
		cpu.pc += uint16(Bytes[ir])
	}
}

func (cpu *Cpu) read(address uint16) byte {
	// switch address {
	// case 0x2002:
	// 	status := cpu.ppu.status
	// 	cpu.ppu.vBlank(false)
	// 	return status
	// }
	return cpu.ram[address]
}

func (cpu *Cpu) write(address uint16, val uint8) {
	switch address {
	// case 0x2000:
	// 	cpu.ppu.controlRegister1(val)
	// case 0x2001:
	// 	cpu.ppu.controlRegister2(val)
	default:
		cpu.ram[address] = val
	}
}

func (cpu *Cpu) reset() {
	// OxFFFC & 0xFFFD contains the intital PC register
	cpu.pc = (uint16(cpu.ram[0xFFFD]) << 8) | uint16(cpu.ram[0xFFFC])
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
		fatal("invalid address")
	}

	var val byte
	if immediate {
		val = cpu.read(cpu.pc + 1)
	} else {
		val = cpu.read(address)
	}

	return val
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
	fmt.Printf("pc: %X, ac: %X , x: %X, y: %X, stack: %X status: %b\n", cpu.pc, cpu.ac, cpu.x, cpu.y, cpu.stack, cpu.getStatus())
}

func Xor(a, b bool) bool {
	return (a || b) && !(a && b)
}

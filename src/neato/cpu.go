package main

import (
	"fmt"
	"os"
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
	"", "izx", "CRASH", "izx", "zp", "zp", "zp", "zp",
	"", "", "", "", "abs", "abs", "abs", "abs",
	"", "izy", "CRASH", "izy", "zpx", "zpx", "zpx", "zpx",
	"", "absy", "", "absy", "absx", "absx", "absx", "absx",
	"", "izx", "CRASH", "izx", "zp", "zp", "zp", "zp",
	"", "", "", "", "abs", "abs", "abs", "abs",
	"", "izy", "CRASH", "izy", "zpx", "zpx", "zpx", "zpx",
	"", "absy", "", "absy", "absx", "absx", "absx", "absx",
	"", "izx", "CRASH", "izx", "zp", "zp", "zp", "zp",
	"", "", "", "", "abs", "abs", "abs", "abs",
	"", "izy", "CRASH", "izy", "zpx", "zpx", "zpx", "zpx",
	"", "absy", "", "absy", "absx", "absx", "absx", "absx",
	"", "izx", "CRASH", "izx", "zp", "zp", "zp", "zp",
	"", "", "", "", "", "abs", "abs", "abs",
	"", "izy", "CRASH", "izy", "zpx", "zpx", "zpx", "zpx",
	"", "absy", "", "absy", "absx", "absx", "absx", "absx",
	"", "izx", "", "izx", "zp", "zp", "zp", "zp",
	"", "", "", "", "abs", "abs", "abs", "abs",
	"", "izy", "CRASH", "izy", "zpx", "zpx", "zpy", "zpy",
	"", "absy", "", "", "absx", "absx", "absy", "absy",
	"", "izx", "im", "izx", "zp", "zp", "zp", "zp",
	"", "", "", "", "abs", "abs", "abs", "abs",
	"", "izy", "CRASH", "izy", "zpx", "zpx", "zpy", "zpy",
	"", "absy", "", "absy", "absx", "absx", "absy", "absy",
	"", "izx", "", "izx", "zp", "zp", "zp", "zp",
	"", "", "", "", "abs", "abs", "abs", "abs",
	"", "izy", "CRASH", "izy", "zpx", "zpx", "zpx", "zpx",
	"", "absy", "", "absy", "absx", "absx", "absx", "absx",
	"", "izx", "", "izx", "zp", "zp", "zp", "zp",
	"", "", "", "", "abs", "abs", "abs", "abs",
	"", "izy", "CRASH", "izy", "zpx", "zpx", "zpx", "zpx",
	"", "absy", "", "absy", "absx", "absx", "absx", "absx"}

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
	fOverflow         uint8
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
	file, err := os.Open(filename)
	if err != nil {
		fatal("file not found: ", filename)
	}
	count, err := file.Read(cpu.ram)
	if err != nil || count != 0xFFF6 {
		fatal("invalid file")
	}
	cpu.pc = 0x1008
	return cpu
}

func (cpu *Cpu) step() {
	cpu.lastPc = cpu.pc

	ir := cpu.read(cpu.pc)
	fmt.Printf("ir 0x%X mode %s cycles %d opcode %s byte %d addr 0x%X\n", ir, AddressingMode[ir], Cycles[ir], Opcodes[ir], Bytes[ir], cpu.pc)

	address := uint16(0)
	immediate := false

	switch AddressingMode[ir] {
	case "abs":
		address = uint16(cpu.read(cpu.pc+2))<<8 | uint16(cpu.read(cpu.pc+1))
	case "zp":
		address = uint16(cpu.read(cpu.pc + 1))

	case "izy":
		base := uint16((cpu.read(cpu.pc + 1)))
		address = uint16(cpu.read(base+1))<<8 | uint16(cpu.read(base))
		address += uint16(cpu.y)

	case "im":
		immediate = true
	case "": // never mind
	default:
		fmt.Println("addressing mode not implemented")
	}

	if address != 0 {
		fmt.Printf("address 0x%X\n", address)
	}

	switch Opcodes[ir] {
	case "CLD":
		cpu.fDecimal = false
	case "CLC":
		cpu.fCarry = false
	case "SEI":
		cpu.fInterruptDisable = true
	case "LDA":
		cpu.ac = cpu.read(address)
		cpu.zeroNeg(cpu.ac)
	case "LDX":
		cpu.x = cpu.val(immediate, address)
		cpu.zeroNeg(cpu.x)
	case "LDY":
		cpu.y = cpu.val(immediate, address)
		cpu.zeroNeg(cpu.x)
	case "CMP":
		temp := cpu.ac - cpu.val(immediate, address)
		cpu.fCarry = temp >= 0
		cpu.zeroNeg(temp)
	case "CPY":
		temp := cpu.y - cpu.val(immediate, address)
		cpu.fCarry = temp >= 0
		cpu.zeroNeg(temp)
	case "ADC":
		cpu.ac = cpu.ac + uint8(cpu.val(immediate, address)) + cpu.carry()
		if uint16(cpu.ac)+uint16(cpu.val(immediate, address))+uint16(cpu.carry()) > 255 {
			cpu.fCarry = true
		}
		cpu.zeroNeg(cpu.ac)
	case "BPL":
		if !cpu.fNegative {
			cpu.pc = uint16(int(cpu.pc) + int(cpu.read(cpu.pc+1)) - 2)
		}
	case "BNE":
		if !cpu.fZero {
			cpu.pc = uint16(int(cpu.pc) + int(cpu.read(cpu.pc+1)) - 2)
		}
	case "BEQ":
		if cpu.fZero {
			cpu.pc = uint16(int(cpu.pc) + int(cpu.read(cpu.pc+1)) - 2)
		}
	case "BCC":
		if cpu.fCarry {
			cpu.pc = uint16(int(cpu.pc) + int(cpu.read(cpu.pc+1)) - 2)
		}
	case "JSR":
		cpu.push(uint8((cpu.pc + 2) >> 8))
		cpu.push(uint8((cpu.pc + 2) & 0xFF))
		cpu.pc = address
	case "JMP":
		fmt.Printf("jmp 0x%X\n", address)
		cpu.pc = address
		// page boundary fix
	case "STX":
		cpu.write(address, cpu.x)
	case "STY":
		cpu.write(address, cpu.y)
	case "STA":
		cpu.write(address, cpu.ac)
	case "DEX":
		cpu.x -= 1
		cpu.zeroNeg(cpu.x)
	case "DEY":
		cpu.y -= 1
		cpu.zeroNeg(cpu.y)
	case "DEC":
		temp := cpu.read(address) - 1
		cpu.zeroNeg(temp)
		cpu.write(address, temp)
	case "TXS":
		cpu.stack = cpu.x
	case "TXA":
		cpu.ac = cpu.x
		cpu.zeroNeg(cpu.ac)
	case "TAX":
		cpu.x = cpu.ac
		cpu.zeroNeg(cpu.x)
	case "TYA":
		cpu.ac = cpu.y
		cpu.zeroNeg(cpu.ac)
	case "EOR":
		cpu.ac = cpu.ac ^ cpu.val(immediate, address)
		cpu.zeroNeg(cpu.ac)
	case "BIT":
		temp := cpu.ac & cpu.val(immediate, address)
		cpu.zeroNeg(temp)
		cpu.fOverflow = (temp >> 6) & 1
	case "PLA":
		cpu.ac = cpu.pull()
		cpu.zeroNeg(cpu.ac)
	default:
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
	cpu.write(0x1000|uint16(cpu.stack), val)
	cpu.stack -= 1
}

func (cpu *Cpu) pull() uint8 {
	cpu.stack += 1
	return cpu.read(0x1000 | uint16(cpu.stack-1))
}

func (cpu *Cpu) val(immediate bool, address uint16) uint8 {
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

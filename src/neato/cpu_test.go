package main

import (
	"fmt"
	"neato/petascii"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func fileDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func newTestCpu(filename string) *Cpu {
	cpu := Cpu{}
	cpu.ram = NewMemory(0xFFFF)
	cpu.printer = petascii.New()
	cpu.respectDecimalMode = true
	cpu.loadFile(filename)
	return &cpu
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

	count, err = file.Read(cpu.ram.store[start:])
	if err != nil || count == 0 {
		fatal("invalid file")
	}

	cpu.ram.Write(0x0002, 0x00)
	cpu.ram.Write(0xA002, 0x00)
	cpu.ram.Write(0xA003, 0x80)
	cpu.ram.Write(0xFFFE, 0x48)
	cpu.ram.Write(0xFFFF, 0xFF)
	cpu.ram.Write(0x01FE, 0xFF)
	cpu.ram.Write(0x01FF, 0x7F)

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

	cpu.ram.Copy(0xFF48, 0xFF5A, []byte{0x48, 0x8A, 0x48, 0x98, 0x48, 0xBA,
		0xBD, 0x04, 0x01, 0x29, 0x10, 0xF0,
		0x03, 0x6C, 0x16, 0x03, 0x6C, 0x14, 0x03})

	cpu.status(0x04)
	cpu.stack = 0xFD
	cpu.pc = 0x0801
}

func (cpu *Cpu) handleTestTraps(t *testing.T) bool {
	switch cpu.pc {
	case 0xFFD2:
		cpu.ram.Write(0x030C, 0)
		cpu.printer.Print(cpu.ac)
		lo := cpu.pull()
		hi := cpu.pull()
		cpu.pc = uint16(hi)<<8 | uint16(lo)
		cpu.pc += 1

	case 0xE16F:
		// each binary are loaded individually.
		return false

	case 0xFFE4:
		cpu.ac = 3
		lo := cpu.pull()
		hi := cpu.pull()
		cpu.pc = uint16(hi)<<8 | uint16(lo)
		cpu.pc += 1

	case 0x8000, 0xA474:
		t.Log("exit trap")
		t.Fail()
	}

	return true
}

func TestCpu(t *testing.T) {
	filepath.Walk(filepath.Join(fileDir(), "cpu_tests", "bin"), func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			runTestFrom(t, path)
		}

		return nil
	})
}

func runTestFrom(t *testing.T, path string) {
	basename := filepath.Base(path)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("", "test failed", basename, r)
		}
	}()

	fmt.Print("\nfile - ", basename)
	cpu := newTestCpu(path)
	for {
		if !cpu.handleTestTraps(t) {
			break
		}

		cpu.Step()
	}
}

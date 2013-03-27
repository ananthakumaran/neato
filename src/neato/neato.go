package main

import (
	"fmt"
	"os"
)

var logInfo = os.Getenv("INFO") == "t"
var logDebug = os.Getenv("DEBUG") == "t"

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("usage neato filename")
		os.Exit(1)
	}

	mapper := LoadRom(args[1])
	ppu := NewPpu(mapper)
	ppu.gui.Init()
	cpu := NewCpu(mapper, ppu)

	for {
		run(cpu)
	}
}

func run(cpu *Cpu) {
	cycles := cpu.Step()
	info(" CYC:%3d SL:%d\n", cpu.ppu.currentScanlineCycle, cpu.ppu.scanline)
	cycles = cycles * 3
	for ; cycles > 0; cycles-- {
		cpu.ppu.Step()
	}
}

func fatal(message ...interface{}) {
	fmt.Println(message...)
	//os.Exit(1)
	panic(message)
}

func debug(format string, message ...interface{}) {
	if logDebug {
		fmt.Printf(format, message...)
	}
}

func info(format string, message ...interface{}) {
	if logInfo {
		fmt.Printf(format, message...)
	}
}

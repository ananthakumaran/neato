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
		fatal("usage neato filename")
	}

	rom := LoadRom(args[1])
	ppu := newPpu(rom)
	cpu := newCpu(rom, ppu)

	for {
		cycles := cpu.step()
		info(" CYC:%3d SL:%d\n", ppu.currentScanlineCycle, ppu.scanline)
		cycles = cycles * 3
		for ; cycles > 0; cycles-- {
			cpu.ppu.step()
		}
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

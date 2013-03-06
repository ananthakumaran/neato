package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fatal("usage neato filename")
	}

	rom := LoadRom(args[1])
	ppu := newPpu(rom)
	cpu := newCpu(rom, ppu)

	cycles := 0
	// brk := 0
	for {
		cycles += cpu.step()
		// if cpu.fBreak {
		// 	ppu.drawScreen()
		// 	brk++
		// 	if brk == 50 {
		// 		break
		// 	}

		// 	cpu.fBreak = false
		// }

		// cycles++

		if cycles%1000 == 0 {
			//debug("\n cycle %d \n", cycles)
		}

		if cycles > 262*113 {
			cpu.ppu.startVblank()
			cycles = 0
		}
	}
}

func fatal(message ...interface{}) {
	fmt.Println(message...)
	os.Exit(1)
}

func debug(format string, message ...interface{}) {
	if os.Getenv("DEBUG") == "t" {
		fmt.Printf(format, message...)
	}
}

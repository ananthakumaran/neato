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

	// rom := LoadRom(args[1])
	// ppu := newPpu(rom)
	// cpu := newCpu(rom, ppu)
	// cpu := testCpu(args[1])
	// for {
	// 	cpu.step()
	// }
}

func fatal(message ...interface{}) {
	fmt.Println(message...)
	os.Exit(1)
}

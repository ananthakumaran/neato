package main

import (
	"fmt"
	"os"
)

var enableGui = os.Getenv("GUI") == "t"

func DrawPixel(x, y int, red, green, blue byte) {
	if !enableGui {
		return
	}
	if x == 0 {
		fmt.Printf("\x1b[0m\n")
	}

	// if x < 250 && y < 150 {
	fmt.Printf("\x1b[48;5;%dm ", 16+(red*36)+(green*6)+blue)
	// }
}

func RefreshScreen() {
	if !enableGui {
		return
	}

	fmt.Printf("%c[2J", 27)
}

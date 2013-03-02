package petascii

import (
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	printer := New()
	for i := uint8(0); i < 0xFF; i++ {
		if i%16 == 0 {
			fmt.Println()
		}
		printer.shift = false
		printer.Print(i)
		printer.shift = true
		printer.Print(i)
		fmt.Print(" ")
	}
}

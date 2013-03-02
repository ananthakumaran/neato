package petascii

import (
	"fmt"
	"strings"
)

type Printer struct {
	shift bool
	text  bool
}

func New() *Printer {
	petascii := Printer{}
	petascii.shift = false
	return &petascii
}

func (p *Printer) Print(code uint8) {
	switch {
	case code <= 0x07:
	case code == 0x08:
		p.shift = false
	case code == 0x09:
		p.shift = true
	case code <= 0x0C:
	case code == 0x0D:
		fmt.Print("\n")
	case code <= 0x13:
	case code == 0x14:
		fmt.Print("\b")
	case code <= 0x1F:
	case code <= 0x40:
		fmt.Printf("%c", code)
	case code <= 0x5A:
		if p.shift {
			fmt.Print(strings.ToLower(fmt.Sprintf("%c", code)))
		} else {
			fmt.Printf("%c", code)
		}
	case code <= 0x60:
		fmt.Print([]string{"[", "£", "]", "↑", "←", "_"}[code-0x5B])
	case code <= 0x7A:
		if p.shift {
			fmt.Print(strings.ToUpper(fmt.Sprintf("%c", code)))
		} else {
			fmt.Print([26]string{
				"♠", "│", "━", "�", "�", "�", "�", "�", "╮",
				"╰", "╯", "�", "╲", "╱", "�", "�", "●", "�",
				"♥", "�", "╭", "╳", "○", "♣", "�", "♦"}[code-0x61])
		}

	case code <= 0x7D:
		fmt.Print([3]string{"┼", "�", "│"}[code-0x7B])
	case code <= 0x7F:
		if p.shift {
			fmt.Print([2]string{"▒", "�"}[code-0x7E])
		} else {
			fmt.Print([2]string{"π", "◥"}[code-0x7E])
		}

	case code <= 0x92:
	case code == 0x93:
		fmt.Printf("\f")
	case code <= 0x9F:
	case code <= 0xA8:
		fmt.Print([9]string{" ", "▌", "▄", "▔", "▁", "▏", "▒", "▕", "�"}[code-0xA0])
	case code == 0xA9:
		if p.shift {
			fmt.Print("�")
		} else {
			fmt.Print("◤")
		}

	case code <= 0xB9:
		fmt.Print([16]string{
			"�", "├", "�", "└", "┐", "▂", "┌", "┴",
			"┬", "┤", "▎", "▍", "�", "�", "�", "▃"}[code-0xAA])

	case code == 0xBA:
		if p.shift {
			fmt.Print("�")
		} else {
			fmt.Print("✓")
		}

	case code <= 0xC0:
		fmt.Print([6]string{"�", "�", "┘", "�", "�", "━"}[code-0xBB])
	case code <= 0xDA:
		if p.shift {
			fmt.Print(strings.ToUpper(fmt.Sprintf("%c", code-96)))
		} else {
			fmt.Print([26]string{
				"♠", "│", "━", "�", "�", "�", "�", "�", "╮",
				"╰", "╯", "�", "╲", "╱", "�", "�", "●", "�",
				"♥", "�", "╭", "╳", "○", "♣", "�", "♦"}[code-0xC1])
		}
	case code <= 0xDD:
		fmt.Print([3]string{"┼", "�", "│"}[code-0xDB])
	case code <= 0xDF:
		if p.shift {
			fmt.Print([2]string{"▒", "�"}[code-0xDE])
		} else {
			fmt.Print([2]string{"π", "◥"}[code-0xDE])
		}

	case code <= 0xE8:
		fmt.Print([9]string{" ", "▌", "▄", "▔", "▁", "▏", "▒", "▕", "�"}[code-0xE0])
	case code == 0xE9:
		if p.shift {
			fmt.Print("�")
		} else {
			fmt.Print("◤")
		}

	case code <= 0xF9:
		fmt.Print([16]string{
			"�", "├", "�", "└", "┐", "▂", "┌", "┴",
			"┬", "┤", "▎", "▍", "�", "�", "�", "▃"}[code-0xEA])
	case code == 0xFA:
		if p.shift {
			fmt.Print("�")
		} else {
			fmt.Print("✓")
		}

	case code <= 0xFE:
		fmt.Print([4]string{"�", "�", "┘", "�"}[code-0xFB])
	case code == 0xFF:
		if p.shift {
			fmt.Print("▒")
		} else {
			fmt.Print("π")
		}

	}
}

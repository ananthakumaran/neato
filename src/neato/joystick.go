package main

import (
	"github.com/go-gl/glfw"
)

type Joystick struct {
	a        int
	b        int
	seleckt  int
	start    int
	up       int
	down     int
	left     int
	right    int
	stickOne int
	stickTwo int
}

func newJoystick() *Joystick {
	joystick := Joystick{}
	glfw.SetKeyCallback(func(key, state int) { joystick.onKey(key, state) })
	return &joystick
}

func (joystick *Joystick) read(address uint16) byte {
	data := 0
	switch address {
	case 0x4016:
		joystick.stickOne++
		switch joystick.stickOne {
		case 1:
			data = joystick.a
		case 2:
			data = joystick.b
		case 3:
			data = joystick.seleckt
		case 4:
			data = joystick.start
		case 5:
			data = joystick.up
		case 6:
			data = joystick.down
		case 7:
			data = joystick.left
		case 8:
			data = joystick.right
		default:
			data = 1
		}
	case 0x4017:
		joystick.stickTwo++
	case 1, 2, 3, 4, 5, 6, 7, 8:
		data = 0
	default:
		data = 1
	}

	joystick.stickOne %= 24
	joystick.stickTwo %= 24

	return byte(data)
}

func (joystick *Joystick) write(address uint16, value byte) {
	switch address {
	case 0x4016:
		joystick.stickOne = 0
	case 0x4017:
		joystick.stickTwo = 0
	}
}

func (joystick *Joystick) onKey(key, state int) {
	switch key {
	case glfw.KeyUp:
		joystick.up = state
	case glfw.KeyDown:
		joystick.down = state
	case glfw.KeyLeft:
		joystick.left = state
	case glfw.KeyRight:
		joystick.right = state
	case glfw.KeyEnter:
		joystick.seleckt = state
	case glfw.KeySpace:
		joystick.start = state
	case 65:
		joystick.a = state
	case 83:
		joystick.b = state
	}
}

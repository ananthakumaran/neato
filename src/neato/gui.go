package main

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"os"
)

var enableGui = os.Getenv("GUI") == "t"

type Gui struct {
	pixels []byte
}

func newGui() *Gui {
	gui := Gui{}
	gui.pixels = make([]byte, SCREEN_WIDTH*SCREEN_HEIGHT*3)

	if err := glfw.Init(); err != nil {
		fatal("can't init glfw", err)
	}

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)

	if err := glfw.OpenWindow(SCREEN_WIDTH*3, SCREEN_HEIGHT*3, 8, 8, 8, 0, 0, 0, glfw.Windowed); err != nil {
		fatal("can't open window", err)
	}

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle("NEato")

	return &gui
}

func (gui *Gui) DrawPixel(x, y int, red, green, blue byte) {
	if !enableGui {
		return
	}

	// opengl starts to draw from the lower left
	base := (SCREEN_HEIGHT - y - 1) * SCREEN_WIDTH * 3
	base += x * 3

	gui.pixels[base] = red
	gui.pixels[base+1] = green
	gui.pixels[base+2] = blue

}

func (gui *Gui) RefreshScreen() {
	if !enableGui {
		return
	}

	gl.PixelZoom(3, 3)
	gl.DrawPixels(SCREEN_WIDTH, SCREEN_HEIGHT, gl.RGB, gl.UNSIGNED_BYTE, gui.pixels)
	glfw.SwapBuffers()
}

package main

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
)

type Gui struct {
	pixels       []byte
	lastMeasured float64
	throttleCnt  int
}

const (
	SCALE                  = 3
	FPS                    = 60.0
	FPS_THROTTLE_FREQUENCY = 15.0
)

func newGui() *Gui {
	gui := Gui{}
	gui.pixels = make([]byte, SCREEN_WIDTH*SCREEN_HEIGHT*3)
	gui.lastMeasured = glfw.Time()

	if err := glfw.Init(); err != nil {
		fatal("can't init glfw", err)
	}

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)

	if err := glfw.OpenWindow(SCREEN_WIDTH*SCALE, SCREEN_HEIGHT*SCALE, 8, 8, 8, 0, 0, 0, glfw.Windowed); err != nil {
		fatal("can't open window", err)
	}

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle("NEato")

	return &gui
}

func (gui *Gui) DrawPixel(x, y int, red, green, blue byte) {
	// opengl starts to draw from the lower left
	base := (SCREEN_HEIGHT - y - 1) * SCREEN_WIDTH * 3
	base += x * 3

	gui.pixels[base] = red
	gui.pixels[base+1] = green
	gui.pixels[base+2] = blue

}

func (gui *Gui) RefreshScreen() {
	gl.PixelZoom(SCALE, SCALE)
	gl.DrawPixels(SCREEN_WIDTH, SCREEN_HEIGHT, gl.RGB, gl.UNSIGNED_BYTE, gui.pixels)
	glfw.SwapBuffers()
	gui.throttle()
}

func (gui *Gui) throttle() {
	gui.throttleCnt++
	if gui.throttleCnt == FPS_THROTTLE_FREQUENCY {
		now := glfw.Time()
		diff := (FPS_THROTTLE_FREQUENCY / FPS) - (now - gui.lastMeasured)

		if diff > 0 {
			if diff > 0.05 {
				glfw.Sleep(diff)
			}
		} else {
			// running slow
		}

		gui.throttleCnt = 0
		gui.lastMeasured = now
	}

}

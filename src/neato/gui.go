package main

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"image"
	"image/color"
)

type Gui struct {
	pixels       []byte
	lastMeasured float64
	throttleCnt  int
	enabled      bool
}

const (
	SCALE                  = 3
	FPS                    = 60.0
	FPS_THROTTLE_FREQUENCY = 15.0
)

func NewGui() *Gui {
	gui := Gui{}
	gui.pixels = make([]byte, SCREEN_WIDTH*SCREEN_HEIGHT*3)
	gui.lastMeasured = glfw.Time()
	gui.enabled = false
	return &gui
}

func (gui *Gui) Init() {
	if err := glfw.Init(); err != nil {
		fatal("can't init glfw", err)
	}

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)

	if err := glfw.OpenWindow(SCREEN_WIDTH*SCALE, SCREEN_HEIGHT*SCALE, 8, 8, 8, 0, 0, 0, glfw.Windowed); err != nil {
		fatal("can't open window", err)
	}

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle("NEato")
	gui.enabled = true
}

func (gui *Gui) Close() {
	glfw.CloseWindow()
	glfw.Terminate()
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
	if !gui.enabled {
		return
	}

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
			gui.lastMeasured = now + diff
			glfw.Sleep(diff)
		} else {
			// running slow
			gui.lastMeasured = now
		}

		gui.throttleCnt = 0

	}

}

func (gui *Gui) TakeScreenShot() image.Image {
	screenshot := image.NewRGBA(image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT))
	for y := 0; y < SCREEN_HEIGHT; y++ {
		for x := 0; x < SCREEN_WIDTH; x++ {
			base := (SCREEN_HEIGHT - y - 1) * SCREEN_WIDTH * 3
			base += x * 3
			screenshot.Set(x, y, color.RGBA{gui.pixels[base], gui.pixels[base+1], gui.pixels[base+2], 255})
		}
	}
	return screenshot
}

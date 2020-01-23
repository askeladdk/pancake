package pancake

import (
	"image"

	"github.com/askeladdk/pancake/graphics"
	"github.com/askeladdk/pancake/input"
)

type Window interface {
	SetCharEventHandler(handler input.CharEventHandler)
	SetKeyEventHandler(handler input.KeyEventHandler)
	SetMouseEventHandler(handler input.MouseEventHandler)
	ShouldClose() bool
	Framebuffer() *graphics.Framebuffer
	SetTitle(string)
	PollEvents()
	SwapBuffers()
}

type Options struct {
	WindowSize image.Point
	Resolution image.Point
	Title      string
	FrameRate  int
}

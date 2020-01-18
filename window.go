package pancake

import (
	"image"

	"github.com/askeladdk/pancake/input"
)

type Window interface {
	SetCharEventHandler(handler input.CharEventHandler)
	SetKeyEventHandler(handler input.KeyEventHandler)
	SetMouseEventHandler(handler input.MouseEventHandler)
	ShouldClose() bool
	Bounds() image.Rectangle
	SetTitle(string)
	PollEvents()
	SwapBuffers()
}

type WindowOptions struct {
	Title string
	Size  image.Point
}

package pancake

import (
	"image"
)

type Window interface {
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

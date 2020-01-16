package pancake

import (
	"image"
)

type Window interface {
	ShouldClose() bool
	Bounds() image.Rectangle
	SetTitle(string)
	SwapBuffers()
}

type WindowOptions struct {
	Title string
	Size  image.Point
}

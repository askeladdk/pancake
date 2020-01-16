package pancake

import (
	"image"
)

type Window interface {
	ShouldClose() bool
	Bounds() image.Rectangle
	Update()
}

type WindowOptions struct {
	Title string
	Size  image.Point
}

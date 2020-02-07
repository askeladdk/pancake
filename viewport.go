package pancake

import (
	"image"
)

func logicalViewport(window, resolution image.Point) image.Rectangle {
	wAspectRatio := float32(window.X) / float32(window.Y)
	rAspectRatio := float32(resolution.X) / float32(resolution.Y)
	var scale int

	if wAspectRatio > rAspectRatio {
		scale = window.Y / resolution.Y
	} else {
		scale = window.X / resolution.X
	}

	logical := resolution.Mul(scale)
	crop := window.Sub(logical)
	border := crop.Div(2)

	return image.Rectangle{
		border,
		border.Add(logical),
	}
}

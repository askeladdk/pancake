package pancake

import (
	"image"
)

func logicalViewport(window, resolution image.Point) image.Rectangle {
	var vw, vh int
	if window.X > window.Y {
		vh = (window.Y / resolution.Y) * resolution.Y
		vw = (vh * resolution.X) / resolution.Y
	} else {
		vw = (window.X / resolution.X) * resolution.X
		vh = (vw * resolution.Y) / resolution.X
	}

	bw := window.X - vw
	bh := window.Y - vh
	borderOffset := image.Point{bw / 2, bh / 2}
	logicalSize := image.Point{window.X - bw, window.Y - bh}
	return image.Rectangle{
		borderOffset,
		borderOffset.Add(logicalSize),
	}
}

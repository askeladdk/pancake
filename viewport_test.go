package pancake

import (
	"image"
	"testing"
)

func TestLogicalViewport(t *testing.T) {
	screen := image.Point{1920, 1080}

	tests := []struct {
		resolution image.Point
		viewport   image.Rectangle
	}{
		{image.Point{640, 360}, image.Rect(0, 0, 1920, 1080)},
		{image.Point{640, 400}, image.Rect(320, 140, 1600, 940)},
		{image.Point{400, 300}, image.Rect(360, 90, 1560, 990)},
	}

	for _, x := range tests {
		logical := logicalViewport(screen, x.resolution)
		if logical != x.viewport {
			t.Fatal(x.resolution)
		}
	}
}

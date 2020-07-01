package tilemap

import (
	"image"

	"github.com/askeladdk/pancake/mathx"
)

// Metrics provides functionality based on the dimensions of a TileMap.
type Metrics struct {
	// Size of cell in pixels. Cells are assumed to be square.
	PixelsPerCell int
	// Width in cells.
	Width int
	// Height in cells.
	Height int
}

func (m Metrics) leptonFromPixel(p int) Lepton {
	return Lepton((p*0x10000 + m.PixelsPerCell/2) / m.PixelsPerCell)
}

func (m Metrics) pixelFromLepton(lp Lepton) int {
	return (int(lp)*m.PixelsPerCell + 0x8000) / 0x10000
}

// Bounds returns the bounds of the grid in pixels.
func (m Metrics) Bounds() image.Rectangle {
	return image.Rect(0, 0, m.Width*m.PixelsPerCell, m.Height*m.PixelsPerCell)
}

// Inside tests if the cell position (x, y) is inside the bounds.
func (m Metrics) Inside(x, y int) bool {
	return x >= 0 && y >= 0 && x < m.Width && y < m.Height
}

// PixelToCoordinate converts a pixel to a Coordinate.
func (m Metrics) PixelToCoordinate(pixel image.Point) Coordinate {
	x := m.leptonFromPixel(pixel.X)
	y := m.leptonFromPixel(pixel.Y)
	return Coord(x, y)
}

// CoordinateToPixel converts a Coordinate to a pixel.
func (m Metrics) CoordinateToPixel(c Coordinate) image.Point {
	x, y := c.xy()
	return image.Pt(m.pixelFromLepton(x), m.pixelFromLepton(y))
}

// AutoTileBitSet computes a bitset by comparing the centre tile
// with all eight neighbouring tiles using a test function.
func (m Metrics) AutoTileBitSet(cx, cy int, testFunc func(x, y int) bool) uint8 {
	var bitset uint8

	for i, offset := range []struct {
		X, Y int
	}{
		// The order determines the meaning of every bit. Do not change!
		{X: +1, Y: -1}, // NE
		{X: +1, Y: +1}, // SE
		{X: -1, Y: +1}, // SW
		{X: -1, Y: -1}, // NW
		{X: +0, Y: -1}, // N
		{X: +1, Y: +0}, // E
		{X: +0, Y: +1}, // S
		{X: -1, Y: +0}, // W
	} {
		x1, y1 := cx+offset.X, cy+offset.Y
		if m.Inside(x1, y1) && testFunc(x1, y1) {
			bitset |= 1 << i
		}
	}

	return bitset
}

// RangeTilesInViewport iterates over all tiles in the viewport
// and returns their positions and modelviews.
func (m Metrics) RangeTilesInViewport(viewport image.Rectangle, fn func(x, y int, modelview mathx.Aff3)) {
	tileSize := mathx.Vec2{float32(m.PixelsPerCell), float32(m.PixelsPerCell)}
	scale := mathx.ScaleAff3(tileSize)
	x0, y0, x1, y1, xofs, yofs := m.drawExtents(viewport)
	pixel := mathx.Vec2{float32(xofs), float32(yofs)}
	xorig := pixel[0]

	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			modelview := scale.Translated(pixel)
			fn(x, y, modelview)
			pixel[0] += tileSize[0]
		}
		pixel[0] = xorig
		pixel[1] += tileSize[1]
	}
}

func (m Metrics) drawExtents(viewport image.Rectangle) (x0, y0, x1, y1, xofs, yofs int) {
	xofs = -(viewport.Min.X % m.PixelsPerCell)
	yofs = -(viewport.Min.Y % m.PixelsPerCell)
	xedge := viewport.Max.X % m.PixelsPerCell
	yedge := viewport.Max.Y % m.PixelsPerCell
	x0 = viewport.Min.X / m.PixelsPerCell
	y0 = viewport.Min.Y / m.PixelsPerCell
	x1 = viewport.Max.X / m.PixelsPerCell
	y1 = viewport.Max.Y / m.PixelsPerCell

	if -xofs+xedge != 0 {
		x1 += 1
	}

	if -yofs+yedge != 0 {
		y1 += 1
	}

	if x1 > m.Width {
		x1 = m.Width
	}

	if y1 > m.Height {
		y1 = m.Height
	}

	return
}

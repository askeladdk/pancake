package tilemap

import (
	"image"

	"github.com/askeladdk/pancake/mathx"
)

// Metrics provides functionality based on the dimensions of a TileMap.
type Metrics struct {
	// Cell size measured in pixels. Cells are always square.
	CellFormat int
	// Width of the map measured in cells.
	CellWidth int
	// Height of the map measured in cells.
	CellHeight int
	// Bounds of the visible area of the map measured in cells.
	CellBounds image.Rectangle
}

func (m Metrics) leptonFromPixel(p int) Lepton {
	return Lepton((p*0x10000 + m.CellFormat/2) / m.CellFormat)
}

func (m Metrics) pixelFromLepton(lp Lepton) int {
	return (int(lp)*m.CellFormat + 0x8000) / 0x10000
}

func (m Metrics) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: m.CellBounds.Min.Mul(m.CellFormat),
		Max: m.CellBounds.Max.Mul(m.CellFormat),
	}
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
		if testFunc(x1, y1) {
			bitset |= 1 << i
		}
	}

	return bitset
}

func (m Metrics) RangeTilesInViewport(viewport image.Rectangle, fn func(x, y int, modelview mathx.Aff3)) {
	tileSize := mathx.Vec2{float32(m.CellFormat), float32(m.CellFormat)}
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
	xofs = -(viewport.Min.X % m.CellFormat)
	yofs = -(viewport.Min.Y % m.CellFormat)
	xedge := viewport.Max.X % m.CellFormat
	yedge := viewport.Max.Y % m.CellFormat
	x0 = viewport.Min.X / m.CellFormat
	y0 = viewport.Min.Y / m.CellFormat
	x1 = viewport.Max.X / m.CellFormat
	y1 = viewport.Max.Y / m.CellFormat

	if -xofs+xedge != 0 {
		x1 += 1
	}

	if -yofs+yedge != 0 {
		y1 += 1
	}

	if x1 > m.CellWidth {
		x1 = m.CellWidth
	}

	if y1 > m.CellHeight {
		y1 = m.CellHeight
	}

	return
}

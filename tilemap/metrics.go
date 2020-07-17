package tilemap

import (
	"image"

	"github.com/askeladdk/pancake/mathx"
)

// Metrics provides functionality based on the dimensions of a TileMap.
type Metrics struct {
	CellFormat CellFormat
	// Width of the map measured in cells.
	CellWidth int
	// Height of the map measured in cells.
	CellHeight int
	// Bounds of the visible area of the map measured in cells.
	CellBounds image.Rectangle
}

func (m Metrics) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: m.CellBounds.Min.Mul(int(m.CellFormat)),
		Max: m.CellBounds.Max.Mul(int(m.CellFormat)),
	}
}

func (m Metrics) CellIndex(coord Coordinate) int {
	x, y := coord.Cell()
	return y*m.CellWidth + x
}

func (m Metrics) RangeTilesInViewport(viewport image.Rectangle, fn func(cell Coordinate, modelview mathx.Aff3)) {
	tileSize := mathx.Vec2{float32(m.CellFormat), float32(m.CellFormat)}
	scale := mathx.ScaleAff3(tileSize)
	x0, y0, x1, y1, xofs, yofs := m.drawExtents(viewport)
	pixel := mathx.Vec2{float32(xofs), float32(yofs)}
	xorig := pixel[0]

	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			modelview := scale.Translated(pixel)
			fn(Cell(x, y), modelview)
			pixel[0] += tileSize[0]
		}
		pixel[0] = xorig
		pixel[1] += tileSize[1]
	}
}

func (m Metrics) drawExtents(viewport image.Rectangle) (x0, y0, x1, y1, xofs, yofs int) {
	xofs = -(viewport.Min.X % int(m.CellFormat))
	yofs = -(viewport.Min.Y % int(m.CellFormat))
	xedge := viewport.Max.X % int(m.CellFormat)
	yedge := viewport.Max.Y % int(m.CellFormat)
	x0 = viewport.Min.X / int(m.CellFormat)
	y0 = viewport.Min.Y / int(m.CellFormat)
	x1 = viewport.Max.X / int(m.CellFormat)
	y1 = viewport.Max.Y / int(m.CellFormat)

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

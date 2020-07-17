package tilemap

import (
	"fmt"
	"image"

	"github.com/askeladdk/pancake/mathx"
)

// Lepton is a pixel-independent measurement of a sub-cell position along a single axis.
type Lepton uint16

// String implements the Stringer interface.
func (l Lepton) String() string {
	return fmt.Sprintf("%02X.%02X", int(l>>8), int(l&0xFF))
}

// Coordinate is a pixel-independent coordinate on a two-dimensional grid.
type Coordinate uint32

const (
	//                                        Y   X
	NorthEast Coordinate = 0xFF000100 // NE (-1, +1)
	SouthEast Coordinate = 0x01000100 // SE (+1, +1)
	SouthWest Coordinate = 0x0100FF00 // SW (+1, -1)
	NorthWest Coordinate = 0xFF00FF00 // NW (-1, -1)
	North     Coordinate = 0xFF000000 // N  (-1,  0)
	East      Coordinate = 0x00000100 // E  ( 0, +1)
	South     Coordinate = 0x01000000 // S  (+1,  0)
	West      Coordinate = 0x0000FF00 // W  ( 0, -1)
)

// Elem returns the coordinate's constituent parts.
func (c Coordinate) Elem() (x, y Lepton) {
	x = Lepton(c & 0xFFFF)
	y = Lepton(c >> 16)
	return
}

// Coord creates a Coordinate from two leptons.
func Coord(x, y Lepton) Coordinate {
	return Coordinate(x) | Coordinate(y)<<16
}

// Cell creates a Coordinate from a cell position.
func Cell(cx, cy int) Coordinate {
	return Coord(Lepton(cx&0xFF)<<8, Lepton(cy&0xFF)<<8)
}

// Cell reports the cell position.
func (c Coordinate) Cell() (cx, cy int) {
	cx = int((c >> 8) & 0xFF)
	cy = int((c >> 24) & 0xFF)
	return
}

// TopLeft returns the Coordinate with the sub-cell position zeroed out.
func (c Coordinate) TopLeft() Coordinate {
	return c & 0xFF00FF00
}

// SubCell returns the Coordinate with the cell position zeroed out.
func (c Coordinate) SubCell() Coordinate {
	return c & 0x00FF00FF
}

// Centered returns the Coordinate at its cell center point.
func (c Coordinate) Centered() Coordinate {
	return c.TopLeft() | 0x00800080
}

// Offset the coordinate by another.
func (c0 Coordinate) Offset(c1 Coordinate) Coordinate {
	x0, y0 := c0.Elem()
	x1, y1 := c1.Elem()
	return Coord(x0+x1, y0+y1)
}

// DistanceVector computes the distance vector between two coordinates.
func (c0 Coordinate) DistanceVector(c1 Coordinate) mathx.Vec2 {
	x0, y0 := c0.Elem()
	x1, y1 := c1.Elem()
	return mathx.Vec2{
		(float32(x1) - float32(x0)) / 0x100,
		(float32(y1) - float32(y0)) / 0x100,
	}
}

// Distance between two coordinates measured in leptons.
func (c0 Coordinate) Distance(c1 Coordinate) Lepton {
	v := c0.DistanceVector(c1)
	return Lepton(0x100 * v.Len())
}

// Move one step in the direction of the target coordinate.
func (c0 Coordinate) Move(c1 Coordinate, step Lepton) Coordinate {
	v := c0.DistanceVector(c1).Unit().Mul(float32(step))
	return c0.Offset(Coord(Lepton(v[0]), Lepton(v[1])))
}

// String implements the Stringer interface.
func (c Coordinate) String() string {
	x, y := c.Elem()
	return fmt.Sprintf("(%v,%v)", x, y)
}

// CellFormat defines the (square) area of a cell in pixels and converts between pixels and coordinates.
type CellFormat int

// PixelToLepton converts a one-dimensional pixel to a lepton.
func (s CellFormat) PixelToLepton(p int) Lepton {
	return Lepton((p*0x100 + int(s/2)) / int(s))
}

// LeptonToPixel converts a one-dimensional lepton to a pixel.
func (s CellFormat) LeptonToPixel(lp Lepton) int {
	return (int(lp)*int(s) + 0x80) / 0x100
}

// Pixel converts a coordinate to a pixel.
func (s CellFormat) Pixel(coord Coordinate) image.Point {
	x, y := coord.Elem()
	return image.Pt(s.LeptonToPixel(x), s.LeptonToPixel(y))
}

// Coordinate converts a pixel to a coordinate.
func (s CellFormat) Coordinate(pixel image.Point) Coordinate {
	x := s.PixelToLepton(pixel.X)
	y := s.PixelToLepton(pixel.Y)
	return Coord(x, y)
}

package tilemap

import (
	"fmt"
)

// Lepton is a pixel-independent measurement of a sub-cell position along a single axis.
type Lepton uint32

// Coordinate is a pixel-independent coordinate on a two-dimensional grid.
//
// A coordinate is made up of two parts:
//
// - A cell position that identifies the position on the TileMap.
//
// - A sub-cell position measured in leptons that identifies the position within the cell.
// This is used for movement and positioning within a cell.
//
// The coordinate system supports a maximum map size of 2^16 x 2^16 cells.
//
// The coordinate system does not specify the tile pixel size.
// Use Metrics to convert to and from pixels.
type Coordinate uint64

func (c Coordinate) xy() (Lepton, Lepton) {
	return Lepton(c & 0xFFFFFFFF), Lepton(c >> 32)
}

// Coord creates a Coordinate from two leptons.
func Coord(x, y Lepton) Coordinate {
	return Coordinate(x) | Coordinate(y)<<32
}

const (
	//                                                      Y   X
	OffsetNorthEast Coordinate = 0xFFFF000000010000 // NE (-1, +1)
	OffsetSouthEast Coordinate = 0x0001000000010000 // SE (+1, +1)
	OffsetSouthWest Coordinate = 0x00010000FFFF0000 // SW (+1, -1)
	OffsetNorthWest Coordinate = 0xFFFF0000FFFF0000 // NW (-1, -1)
	OffsetNorth     Coordinate = 0xFFFF000000000000 // N  (-1,  0)
	OffsetEast      Coordinate = 0x0000000000010000 // E  ( 0, +1)
	OffsetSouth     Coordinate = 0x0001000000000000 // S  (+1,  0)
	OffsetWest      Coordinate = 0x00000000FFFF0000 // W  ( 0, -1)
)

// CellXY returns the x and y position of the referenced cell.
func (c Coordinate) CellXY() (x, y int) {
	x = int((c >> 16) & 0xFFFF)
	y = int((c >> 48) & 0xFFFF)
	return
}

// Cell returns the Coordinate with the sub-cell position zeroed out.
func (c Coordinate) Cell() Coordinate {
	return c & 0xFFFF0000FFFF0000
}

// SubCell returns the Coordinate with the cell position zeroed out.
func (c Coordinate) SubCell() Coordinate {
	return c & 0x0000FFFF0000FFFF
}

// Centered returns the Coordinate at its cell center point.
func (c Coordinate) Centered() Coordinate {
	return c.Cell() | 0x0000800000008000
}

func (c0 Coordinate) Add(c1 Coordinate) Coordinate {
	x0, y0 := c0.xy()
	x1, y1 := c1.xy()
	return Coord(x0+x1, y0+y1)
}

func (c0 Coordinate) Sub(c1 Coordinate) Coordinate {
	x0, y0 := c0.xy()
	x1, y1 := c1.xy()
	return Coord(x0-x1, y0-y1)
}

func (c Coordinate) String() string {
	x, y := c.xy()
	return fmt.Sprintf("(%d.%04X,%d.%04X)", int(x>>16), int(x&0xFFFF), int(y>>16), int(y&0xFFFF))
}

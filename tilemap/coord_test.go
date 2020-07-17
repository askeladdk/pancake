package tilemap

import (
	"fmt"
	"image"
	"testing"
)

func Test_Coordinate_Pixel(t *testing.T) {
	s := CellFormat(16)
	p := image.Pt(100, 125)
	c := s.Coordinate(p)
	if s.Pixel(c) != p {
		t.Fatal()
	}
}

func Test_Cell(t *testing.T) {
	c := Cell(12, 13).Centered()
	if x, y := c.Cell(); x != 12 || y != 13 {
		t.Fatal(x, y)
	} else if c.SubCell() != 0x00800080 {
		t.Fatal()
	}
}

func Test_Compass(t *testing.T) {
	c := Cell(1, 1)
	for _, test := range []struct {
		ofs  Coordinate
		x, y int
	}{
		{North, 1, 0},
		{East, 2, 1},
		{South, 1, 2},
		{West, 0, 1},
		{NorthEast, 2, 0},
		{NorthWest, 0, 0},
		{SouthEast, 2, 2},
		{SouthWest, 0, 2},
	} {
		if x, y := c.Offset(test.ofs).Cell(); x != test.x || y != test.y {
			t.Fatal(test.ofs, test.x, test.y)
		}
	}
}

func Test_Wrap(t *testing.T) {
	c := Coordinate(0x0100FF00)
	if x := c.Offset(SouthEast); x != 0x02000000 {
		t.Fatal(x)
	}
}

func Test_Move(t *testing.T) {
	s := CellFormat(16)
	c0 := Cell(0, 0)
	c1 := Cell(2, 1)
	step := Lepton(0x0090)
	dist := c0.Distance(c1)
	for dist > step {
		c0 = c0.Move(c1, step)
		dist = c0.Distance(c1)
		fmt.Println(c0, s.Pixel(c0))
	}
	c0 = c1
	fmt.Println(c0, s.Pixel(c0))
}

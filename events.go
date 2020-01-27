package pancake

import (
	"image"

	"github.com/askeladdk/pancake/input"
)

type CharEvent struct {
	Char rune
}

type KeyEvent struct {
	Key      input.Key
	Flags    input.Flags
	Scancode int
}

type MouseEvent struct {
	Mouse    input.Mouse
	Flags    input.Flags
	Position image.Point
}

type FrameEvent struct {
	DeltaTime float64
}

type DrawEvent struct {
	Alpha float64
}

package pancake

import (
	"image"

	"github.com/askeladdk/pancake/input"
)

type QuitEvent struct{}

type CharEvent struct {
	Char rune
}

type KeyEvent struct {
	Key      input.Key
	Flags    input.Flags
	Scancode int
}

type MouseEvent struct {
	Button   input.MouseButton
	Flags    input.Flags
	Position image.Point
}

type MouseMoveEvent struct {
	Position image.Point
}

type FrameEvent struct {
	DeltaTime float64
}

type DrawEvent struct {
	Alpha float64
}

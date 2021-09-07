package pancake

import (
	"image"
)

type CloseEvent struct{}

type CharEvent struct {
	Char rune
}

type KeyEvent struct {
	Modifiers
	Key      Key
	Scancode int
}

type MouseEvent struct {
	Modifiers
	Button   MouseButton
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

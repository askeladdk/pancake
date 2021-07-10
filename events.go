package pancake

import (
	"image"
)

type CloseEvent struct{}

type CharEvent struct {
	Char rune
}

type KeyEvent struct {
	Key       Key
	Modifiers Modifiers
	Scancode  int
}

type MouseEvent struct {
	Button    MouseButton
	Modifiers Modifiers
	Position  image.Point
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

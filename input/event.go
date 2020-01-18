package input

import "image"

type KeyEvent struct {
	Key      Key
	Flags    Flags
	Scancode int
}

type MouseEvent struct {
	Mouse    Mouse
	Flags    Flags
	Position image.Point
}

type CharEventHandler interface {
	CharEvent(rune)
}

type KeyEventHandler interface {
	KeyEvent(KeyEvent)
}

type MouseEventHandler interface {
	MouseEvent(MouseEvent)
}

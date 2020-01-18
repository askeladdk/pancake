package input

import (
	"fmt"
)

type Flags uint

const (
	Pressed Flags = 1 << iota
	Released
	Repeated
	Shift
	Control
	Alt
	Super
)

func (f Flags) String() string {
	return fmt.Sprintf("%08b", f)
}

func (flags Flags) Pressed() bool {
	return flags&Pressed != 0
}

func (flags Flags) Released() bool {
	return flags&Released != 0
}

func (flags Flags) Repeated() bool {
	return flags&Repeated != 0
}

func (flags Flags) Shift() bool {
	return flags&Shift != 0
}

func (flags Flags) Control() bool {
	return flags&Control != 0
}

func (flags Flags) Alt() bool {
	return flags&Alt != 0
}

func (flags Flags) Super() bool {
	return flags&Super != 0
}

type Mouse int

const (
	MouseButton0 Mouse = 0
	MouseButton1 Mouse = 1
	MouseButton2 Mouse = 2
	MouseButton3 Mouse = 3
	MouseButton4 Mouse = 4
	MouseButton5 Mouse = 5
	MouseButton6 Mouse = 6
	MouseButton7 Mouse = 7
	MouseMove    Mouse = 8
)

type Key int

const (
	Key0            Key = 48
	Key1            Key = 49
	Key2            Key = 50
	Key3            Key = 51
	Key4            Key = 52
	Key5            Key = 53
	Key6            Key = 54
	Key7            Key = 55
	Key8            Key = 56
	Key9            Key = 57
	KeyA            Key = 65
	KeyApostrophe   Key = 39
	KeyB            Key = 66
	KeyBackslash    Key = 92
	KeyBackspace    Key = 259
	KeyC            Key = 67
	KeyCapsLock     Key = 280
	KeyComma        Key = 44
	KeyD            Key = 68
	KeyDelete       Key = 261
	KeyDown         Key = 264
	KeyE            Key = 69
	KeyEnd          Key = 269
	KeyEnter        Key = 257
	KeyEqual        Key = 61
	KeyEscape       Key = 256
	KeyF            Key = 70
	KeyF1           Key = 290
	KeyF10          Key = 299
	KeyF11          Key = 300
	KeyF12          Key = 301
	KeyF2           Key = 291
	KeyF3           Key = 292
	KeyF4           Key = 293
	KeyF5           Key = 294
	KeyF6           Key = 295
	KeyF7           Key = 296
	KeyF8           Key = 297
	KeyF9           Key = 298
	KeyG            Key = 71
	KeyGraveAccent  Key = 96
	KeyH            Key = 72
	KeyHome         Key = 268
	KeyI            Key = 73
	KeyInsert       Key = 260
	KeyJ            Key = 74
	KeyK            Key = 75
	KeyKP0          Key = 320
	KeyKP1          Key = 321
	KeyKP2          Key = 322
	KeyKP3          Key = 323
	KeyKP4          Key = 324
	KeyKP5          Key = 325
	KeyKP6          Key = 326
	KeyKP7          Key = 327
	KeyKP8          Key = 328
	KeyKP9          Key = 329
	KeyKPAdd        Key = 334
	KeyKPDecimal    Key = 330
	KeyKPDivide     Key = 331
	KeyKPEnter      Key = 335
	KeyKPEqual      Key = 336
	KeyKPMultiply   Key = 332
	KeyKPSubtract   Key = 333
	KeyL            Key = 76
	KeyLast         Key = 348
	KeyLeft         Key = 263
	KeyLeftAlt      Key = 342
	KeyLeftBracket  Key = 91
	KeyLeftControl  Key = 341
	KeyLeftShift    Key = 340
	KeyLeftSuper    Key = 343
	KeyM            Key = 77
	KeyMenu         Key = 348
	KeyMinus        Key = 45
	KeyN            Key = 78
	KeyNumLock      Key = 282
	KeyO            Key = 79
	KeyP            Key = 80
	KeyPageDown     Key = 267
	KeyPageUp       Key = 266
	KeyPause        Key = 284
	KeyPeriod       Key = 46
	KeyPrintScreen  Key = 283
	KeyQ            Key = 81
	KeyR            Key = 82
	KeyRight        Key = 262
	KeyRightAlt     Key = 346
	KeyRightBracket Key = 93
	KeyRightControl Key = 345
	KeyRightShift   Key = 344
	KeyRightSuper   Key = 347
	KeyS            Key = 83
	KeyScrollLock   Key = 281
	KeySemicolon    Key = 59
	KeySlash        Key = 47
	KeySpace        Key = 32
	KeyT            Key = 84
	KeyTab          Key = 258
	KeyU            Key = 85
	KeyUnknown      Key = -1
	KeyUp           Key = 265
	KeyV            Key = 86
	KeyW            Key = 87
	KeyWorld1       Key = 161
	KeyWorld2       Key = 162
	KeyX            Key = 88
	KeyY            Key = 89
	KeyZ            Key = 90
)

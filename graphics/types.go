package graphics

import (
	"errors"
	"fmt"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

type Filter uint32

const (
	FilterLinear Filter = iota
	FilterNearest
)

func (this Filter) param() gl.Enum {
	switch this {
	case FilterLinear:
		return gl.LINEAR
	case FilterNearest:
		return gl.NEAREST
	default:
		panic(errors.New("invalid filter"))
	}
}

type ColorFormat uint32

const (
	ColorFormatRGBA ColorFormat = iota
	ColorFormatRGB
	ColorFormatIndexed
)

func (this ColorFormat) format() gl.Enum {
	switch this {
	case ColorFormatRGBA:
		return gl.RGBA
	case ColorFormatRGB:
		return gl.RGB
	case ColorFormatIndexed:
		return gl.RED
	default:
		panic(errors.New("invalid color mode"))
	}
}

func (this ColorFormat) internalFormat() gl.Enum {
	switch this {
	case ColorFormatRGBA:
		return gl.RGBA
	case ColorFormatRGB:
		return gl.RGB
	case ColorFormatIndexed:
		return gl.R8
	default:
		panic(errors.New("invalid color mode"))
	}
}

func (this ColorFormat) pixelSize() int {
	switch this {
	case ColorFormatRGBA:
		return 4
	case ColorFormatRGB:
		return 4
	case ColorFormatIndexed:
		return 1
	default:
		panic(errors.New("invalid color mode"))
	}
}

type AttribType uint32

const (
	Float AttribType = iota
	Vec2
	Vec3
	Vec4
	Mat2
	Mat23
	Mat24
	Mat3
	Mat32
	Mat34
	Mat4
	Mat42
	Mat43
)

func (this AttribType) components() int {
	switch this {
	case Float:
		return 1
	case Vec2:
		return 2
	case Vec3:
		return 3
	case Vec4:
		return 4
	case Mat2:
		return 2
	case Mat23:
		return 2
	case Mat24:
		return 2
	case Mat3:
		return 3
	case Mat32:
		return 3
	case Mat34:
		return 3
	case Mat4:
		return 4
	case Mat42:
		return 4
	case Mat43:
		return 4
	default:
		panic(fmt.Errorf("invalid attribute type"))
	}
}

func (this AttribType) repeat() int {
	switch this {
	case Mat2:
		return 2
	case Mat23:
		return 3
	case Mat24:
		return 4
	case Mat3:
		return 3
	case Mat32:
		return 2
	case Mat34:
		return 4
	case Mat4:
		return 4
	case Mat42:
		return 2
	case Mat43:
		return 3
	default:
		return 1
	}
}

func (this AttribType) stride() int {
	return 4 * this.components()
}

type AttribFormat []AttribType

func (this AttribFormat) stride() int {
	var stride int
	for _, a := range this {
		stride += a.stride() * a.repeat()
	}
	return stride
}

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

func (filter Filter) param() gl.Enum {
	switch filter {
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

func (format ColorFormat) format() gl.Enum {
	switch format {
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

func (format ColorFormat) internalFormat() gl.Enum {
	switch format {
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

func (format ColorFormat) pixelSize() int {
	switch format {
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
	Float32 AttribType = iota
	Float64
	Vec2
	Vec3
	Vec4
	Mat3
	Mat4
	Byte4
)

func (atype AttribType) components() int {
	switch atype {
	case Float32:
		return 1
	case Float64:
		return 1
	case Vec2:
		return 2
	case Vec3:
		return 3
	case Vec4:
		return 4
	case Mat3:
		return 3
	case Mat4:
		return 4
	case Byte4:
		return 4
	default:
		panic(fmt.Errorf("invalid attribute type"))
	}
}

func (atype AttribType) repeat() int {
	switch atype {
	case Mat3:
		return 3
	case Mat4:
		return 4
	default:
		return 1
	}
}

func (atype AttribType) bytes() int {
	switch atype {
	case Byte4:
		return 1
	case Float32:
		return 4
	default:
		return 8
	}
}

func (atype AttribType) xtype() gl.Enum {
	switch atype {
	case Byte4:
		return gl.UNSIGNED_BYTE
	case Float32:
		return gl.FLOAT
	default:
		return gl.DOUBLE
	}
}

func (atype AttribType) normalised() bool {
	switch atype {
	case Byte4:
		return true
	default:
		return false
	}
}

func (atype AttribType) stride() int {
	return atype.bytes() * atype.components()
}

type AttribFormat []AttribType

func (aformat AttribFormat) stride() int {
	var stride int
	for _, a := range aformat {
		stride += a.stride() * a.repeat()
	}
	return stride
}

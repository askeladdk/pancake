package graphics

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Filter uint32

const (
	FilterLinear Filter = iota
	FilterNearest
	FilterLinearMipmap
	FilterNearestNearest
)

func (this Filter) minparam() int32 {
	switch this {
	case FilterLinear:
		return gl.LINEAR
	case FilterNearest:
		return gl.NEAREST
	case FilterLinearMipmap:
		return gl.LINEAR_MIPMAP_LINEAR
	case FilterNearestNearest:
		return gl.NEAREST_MIPMAP_NEAREST
	default:
		panic(errors.New("invalid filter"))
	}
}

func (this Filter) maxparam() int32 {
	switch this {
	case FilterLinear:
		return gl.LINEAR
	case FilterNearest:
		return gl.NEAREST
	case FilterLinearMipmap:
		return gl.LINEAR
	case FilterNearestNearest:
		return gl.NEAREST
	default:
		panic(errors.New("invalid filter"))
	}
}

func (this Filter) mipmap() bool {
	switch this {
	case FilterLinear:
		return false
	case FilterNearest:
		return false
	case FilterLinearMipmap:
		return true
	case FilterNearestNearest:
		return true
	default:
		panic(errors.New("invalid filter"))
	}
}

type ColorFormat uint32

const (
	ColorFormatRGBA ColorFormat = iota
	ColorFormatIndexed
	colorFormatDepthStencil
)

func (this ColorFormat) format() uint32 {
	switch this {
	case ColorFormatRGBA:
		return gl.RGBA
	case ColorFormatIndexed:
		return gl.RED
	case colorFormatDepthStencil:
		return gl.DEPTH_STENCIL
	default:
		panic(errors.New("invalid color mode"))
	}
}

func (this ColorFormat) internalFormat() int32 {
	switch this {
	case ColorFormatRGBA:
		return gl.RGBA
	case ColorFormatIndexed:
		return gl.R8
	case colorFormatDepthStencil:
		return gl.DEPTH24_STENCIL8
	default:
		panic(errors.New("invalid color mode"))
	}
}

func (this ColorFormat) xtype() uint32 {
	switch this {
	case ColorFormatRGBA:
		return gl.UNSIGNED_BYTE
	case ColorFormatIndexed:
		return gl.UNSIGNED_BYTE
	case colorFormatDepthStencil:
		return gl.UNSIGNED_INT_24_8
	default:
		panic(errors.New("invalid color mode"))
	}
}

func (this ColorFormat) pixelSize() int {
	switch this {
	case ColorFormatRGBA:
		return 4
	case ColorFormatIndexed:
		return 1
	case colorFormatDepthStencil:
		return 4
	default:
		panic(errors.New("invalid color mode"))
	}
}

type AttrType uint32

const (
	Float AttrType = iota
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

func (this AttrType) components() int32 {
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

func (this AttrType) repeat() int32 {
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

func (this AttrType) stride() int32 {
	return 4 * this.components()
}

type AttrFormat []AttrType

func (this AttrFormat) stride() int32 {
	var stride int32
	for _, a := range this {
		stride += a.stride() * a.repeat()
	}
	return stride
}

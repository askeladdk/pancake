package pancake

import (
	"errors"
	"fmt"

	gl "github.com/askeladdk/pancake/opengl"
)

type TextureFilter uint32

const (
	FilterLinear TextureFilter = iota
	FilterNearest
)

func (filter TextureFilter) param() gl.Enum {
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

type Attrib uint32

const (
	AttribFloat32 Attrib = iota
	AttribFloat64
	AttribVec2
	AttribVec3
	AttribVec4
	AttribMat3
	AttribMat4
	AttribByte4
)

func (atype Attrib) components() int {
	switch atype {
	case AttribFloat32:
		return 1
	case AttribFloat64:
		return 1
	case AttribVec2:
		return 2
	case AttribVec3:
		return 3
	case AttribVec4:
		return 4
	case AttribMat3:
		return 3
	case AttribMat4:
		return 4
	case AttribByte4:
		return 4
	default:
		panic(fmt.Errorf("invalid attribute type"))
	}
}

func (atype Attrib) repeat() int {
	switch atype {
	case AttribMat3:
		return 3
	case AttribMat4:
		return 4
	default:
		return 1
	}
}

func (atype Attrib) bytes() int {
	switch atype {
	case AttribByte4:
		return 1
	case AttribFloat32:
		return 4
	default:
		return 8
	}
}

func (atype Attrib) xtype() gl.Enum {
	switch atype {
	case AttribByte4:
		return gl.UNSIGNED_BYTE
	case AttribFloat32:
		return gl.FLOAT
	default:
		return gl.DOUBLE
	}
}

func (atype Attrib) normalised() bool {
	switch atype {
	case AttribByte4:
		return true
	default:
		return false
	}
}

func (atype Attrib) stride() int {
	return atype.bytes() * atype.components()
}

type AttribFormat []Attrib

func (aformat AttribFormat) stride() int {
	var stride int
	for _, a := range aformat {
		stride += a.stride() * a.repeat()
	}
	return stride
}

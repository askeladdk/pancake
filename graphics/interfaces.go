package graphics

import (
	"image"
)

type Texture interface {
	Begin()
	End()
	BeginAt(unit int)
	EndAt(unit int)
	Size() image.Point
	ColorFormat() ColorFormat
	Filter() Filter
	SetFilter(filter Filter)
	SetPixels(pixels []byte)
	Pixels(pixels []byte) []byte
}

type Shader interface {
	Begin()
	End()
	SetUniform(name string, value interface{}) bool
}

type Frame interface {
	Begin()
	End()
	Texture() Texture
	Blit(dst Frame, dr, sr image.Rectangle, filter Filter)
}

type Buffer interface {
	Begin()
	End()
	Len() int
	SetData(i, j int, data interface{})
}

type VertexSlice interface {
	Begin()
	End()
	Format() AttrFormat
	Len() int
	SetData(data interface{})
	Slice(i, j int) VertexSlice
	Draw()
}

// type DrawArray interface {
// 	Begin()
// 	End()
// 	Slice(b, n int) Buffer
// 	Draw(i, j int) error
// }

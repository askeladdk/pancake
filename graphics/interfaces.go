package graphics

import (
	"image"
)

// var (
// 	Quad = []float32{
// 		-.5, -.5, 0, 0,
// 		-.5, +.5, 0, 1,
// 		+.5, -.5, 1, 0,
// 		-.5, +.5, 0, 1,
// 		+.5, -.5, 1, 0,
// 		+.5, +.5, 1, 1,
// 	}

// 	QuadDescriptor = DataDescriptor{
// 		Count: 6,
// 		Format: AttrFormat{
// 			Vec2,
// 			Vec2,
// 		},
// 	}
// )

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

// type Buffer interface {
// 	Begin()
// 	End()
// 	Format() AttrFormat
// 	Len() int
// 	SetData(i, j int, data []float32)
// }

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

type Driver interface {
	Version() string
	Flush()
	SetViewport(bounds image.Rectangle)
	Clear(r, g, b, a float32)
	NewTexture(size image.Point, filter Filter, format ColorFormat, pixels []byte) (Texture, error)
	NewShader(vshader, fshader string) (Shader, error)
	NewFrame(size image.Point, filter Filter, depthStencil bool) (Frame, error)
	ScreenFrame() Frame
	NewVertexSlice(format AttrFormat, len int, data interface{}) (VertexSlice, error)
	// NewBuffer(format AttrFormat, len int, data []float32) (Buffer, error)
}

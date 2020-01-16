package graphics

import (
	"image"
	"runtime"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

type renderbuffer struct {
	id gl.Renderbuffer
}

func (rbo *renderbuffer) delete() {
	gl.DeleteRenderbuffer(rbo.id)
}

func newRenderbuffer(size image.Point, internalFormat gl.Enum, samples int) *renderbuffer {
	rbo := &renderbuffer{
		id: gl.CreateRenderbuffer(),
	}

	runtime.SetFinalizer(rbo, (*renderbuffer).delete)

	gl.BindRenderbuffer(rbo.id)
	defer gl.BindRenderbuffer(0)
	gl.RenderbufferStorage(internalFormat, size.X, size.Y, samples)

	return rbo
}

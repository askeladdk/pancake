package graphics

import (
	"errors"
	"image"
	"runtime"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

var fbobinder = newBinder(func(id uint32) {
	gl.BindFramebuffer(gl.Framebuffer(id))
})

type Framebuffer struct {
	color   *Texture
	depth   *renderbuffer
	stencil *renderbuffer
	id      gl.Framebuffer
}

func (fbo *Framebuffer) Begin() {
	fbobinder.bind(uint32(fbo.id))
}

func (fbo *Framebuffer) End() {
	fbobinder.unbind()
}

func (fbo *Framebuffer) Texture() *Texture {
	return fbo.color
}

func (fbo *Framebuffer) Blit(dst *Framebuffer, dr, sr image.Rectangle, filter Filter) {
	gl.BlitNamedFramebuffer(
		fbo.id, dst.id,
		sr.Min.X, sr.Min.Y, sr.Max.X, sr.Max.Y,
		dr.Min.X, dr.Min.Y, dr.Max.X, dr.Max.Y,
		gl.COLOR_BUFFER_BIT, filter.param(),
	)
	panicError()
}

func (fbo *Framebuffer) delete() {
	gl.DeleteFramebuffer(fbo.id)
}

func NewFramebuffer(size image.Point, filter Filter, depth, stencil bool) (*Framebuffer, error) {
	fbo := &Framebuffer{
		id:    gl.CreateFramebuffer(),
		color: NewTexture(size, filter, ColorFormatRGBA, nil),
	}

	runtime.SetFinalizer(fbo, (*Framebuffer).delete)

	fbo.Begin()
	defer fbo.End()

	gl.FramebufferTexture2D(gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, fbo.color.id, 0)

	if depth {
		fbo.depth = newRenderbuffer(size, gl.DEPTH_COMPONENT16, 0)
		gl.FramebufferRenderbuffer(gl.DEPTH_ATTACHMENT, fbo.depth.id)
	}

	if stencil {
		fbo.stencil = newRenderbuffer(size, gl.STENCIL_INDEX8, 0)
		gl.FramebufferRenderbuffer(gl.STENCIL_ATTACHMENT, fbo.stencil.id)
	}

	if code := gl.CheckFramebufferStatus(); code != gl.FRAMEBUFFER_COMPLETE {
		return nil, errors.New(errorToString(code))
	} else {
		return fbo, nil
	}
}

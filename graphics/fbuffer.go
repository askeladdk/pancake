package graphics

import (
	"errors"
	"image"
	"runtime"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

var fbobinder = newBinder(func(id uint32) {
	gl.BindFramebuffer(gl.FRAMEBUFFER, gl.Framebuffer(id))
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

func (fbo *Framebuffer) Bounds() image.Rectangle {
	return fbo.color.Bounds()
}

func (fbo *Framebuffer) Color() *Texture {
	return fbo.color
}

func (dst *Framebuffer) BlitFrom(src *Framebuffer, dr, sr image.Rectangle, mask gl.Enum, filter Filter) {
	gl.BlitNamedFramebuffer(src.id, dst.id, sr, dr, mask, filter.param())
}

func (fbo *Framebuffer) delete() {
	gl.DeleteFramebuffer(fbo.id)
}

func NewFramebufferFromTexture(color *Texture, depthStencil bool) (*Framebuffer, error) {
	if color.ColorFormat() != ColorFormatRGB && color.ColorFormat() != ColorFormatRGBA {
		return nil, errors.New("color texture must be in RGB(A) color format")
	}

	fbo := &Framebuffer{
		id:    gl.CreateFramebuffer(),
		color: color,
	}

	runtime.SetFinalizer(fbo, (*Framebuffer).delete)

	fbo.Begin()
	defer fbo.End()

	gl.FramebufferTexture2D(gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, fbo.color.id, 0)

	if depthStencil {
		depthStencilBuffer := newRenderbuffer(color.size, gl.DEPTH24_STENCIL8, 0)
		gl.FramebufferRenderbuffer(gl.DEPTH_STENCIL_ATTACHMENT, depthStencilBuffer.id)
		fbo.depth = depthStencilBuffer
		fbo.stencil = depthStencilBuffer
	}

	if code := gl.CheckFramebufferStatus(); code != gl.FRAMEBUFFER_COMPLETE {
		return nil, errors.New(errorToString(code))
	} else {
		return fbo, nil
	}
}

func NewFramebuffer(size image.Point, filter Filter, depthStencil bool) (*Framebuffer, error) {
	color := NewTexture(size, filter, ColorFormatRGBA, nil)
	return NewFramebufferFromTexture(color, depthStencil)
}

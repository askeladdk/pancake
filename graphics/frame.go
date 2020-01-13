package graphics

import (
	"errors"
	"image"
	"runtime"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

var fboBinder = newBinder(func(ref uint32) {
	gl.BindFramebuffer(gl.FRAMEBUFFER, gl.Framebuffer(ref))
})

type frame struct {
	color *texture
	depth *texture
	ref   gl.Framebuffer
}

func (this *frame) Begin() {
	fboBinder.bind(uint32(this.ref))
}

func (this *frame) End() {
	fboBinder.unbind()
}

func (this *frame) Texture() Texture {
	return this.color
}

func (this *frame) Blit(dst Frame, dr, sr image.Rectangle, filter Filter) {
	gl.BlitNamedFramebuffer(
		this.ref, dst.(*frame).ref,
		sr.Min.X, sr.Min.Y, sr.Max.X, sr.Max.Y,
		dr.Min.X, dr.Min.Y, dr.Max.X, dr.Max.Y,
		gl.COLOR_BUFFER_BIT, filter.maxparam(),
	)
	panicError()
}

func (this *frame) delete() {
	gl.DeleteFramebuffer(this.ref)
}

func NewFrame(size image.Point, filter Filter, depthStencil bool) (*frame, error) {
	if color, err := NewTexture(size, filter, ColorFormatRGBA, nil); err != nil {
		return nil, err
	} else {
		fbo := &frame{
			ref:   gl.CreateFramebuffer(),
			color: color,
		}

		runtime.SetFinalizer(fbo, (*frame).delete)

		fbo.Begin()
		defer fbo.End()

		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, fbo.color.ref, 0)

		if depthStencil {
			if depth, err := NewTexture(size, FilterLinear, colorFormatDepthStencil, nil); err != nil {
				return nil, err
			} else {
				fbo.depth = depth
				gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.TEXTURE_2D, fbo.depth.ref, 0)
			}
		}

		if code := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); code != gl.FRAMEBUFFER_COMPLETE {
			return nil, errors.New(errorToString(code))
		} else {
			return fbo, nil
		}
	}
}

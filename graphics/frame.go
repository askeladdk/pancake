package graphics

import (
	"errors"
	"image"
	"runtime"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var fboBinder = newBinder(func(ref uint32) {
	gl.BindFramebuffer(gl.FRAMEBUFFER, ref)
})

type frame struct {
	color *texture
	depth *texture
	ref   uint32
}

func (this *frame) Begin() {
	fboBinder.bind(this.ref)
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
		int32(sr.Min.X), int32(sr.Min.Y), int32(sr.Max.X), int32(sr.Max.Y),
		int32(dr.Min.X), int32(dr.Min.Y), int32(dr.Max.X), int32(dr.Max.Y),
		gl.COLOR_BUFFER_BIT, uint32(filter.maxparam()),
	)
	panicError()
}

func (this *frame) delete() {
	mainthread.CallNonBlock(func() {
		gl.DeleteFramebuffers(1, &this.ref)
	})
}

func newFrame(size image.Point, filter Filter, depthStencil bool) (*frame, error) {
	if color, err := newTexture(size, filter, ColorFormatRGBA, nil); err != nil {
		return nil, err
	} else {
		fbo := &frame{
			color: color,
		}
		gl.GenFramebuffers(1, &fbo.ref)
		runtime.SetFinalizer(fbo, (*frame).delete)

		fbo.Begin()
		defer fbo.End()

		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, fbo.color.ref, 0)

		if depthStencil {
			if depth, err := newTexture(size, FilterLinear, colorFormatDepthStencil, nil); err != nil {
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

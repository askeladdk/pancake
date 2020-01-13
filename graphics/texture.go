package graphics

import (
	"errors"
	"image"
	"runtime"
	"unsafe"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	TextureUnitsCount = 48 // Minimum number of texture units according to the spec.
)

var texbinders = func(n uint32) []*binder {
	var binders []*binder
	for i := uint32(0); i < n; i++ {
		x := i // gotcha
		binders = append(binders, newBinder(func(ref uint32) {
			gl.ActiveTexture(gl.TEXTURE0 + x)
			gl.BindTexture(gl.TEXTURE_2D, ref)
		}))
	}
	return binders
}(TextureUnitsCount)

type texture struct {
	size   image.Point
	format ColorFormat
	filter Filter
	ref    uint32
	mipmap bool
}

func (this *texture) BeginAt(unit int) {
	texbinders[unit].bind(this.ref)
}

func (this *texture) EndAt(unit int) {
	texbinders[unit].unbind()
}

func (this *texture) Begin() {
	this.BeginAt(0)
}

func (this *texture) End() {
	this.EndAt(0)
}

func (this *texture) Size() image.Point {
	return this.size
}

func (this *texture) ColorFormat() ColorFormat {
	return this.format
}

func (this *texture) Filter() Filter {
	return this.filter
}

func (this *texture) SetFilter(filter Filter) {
	if filter.mipmap() && !this.mipmap {
		gl.GenerateMipmap(gl.TEXTURE_2D)
		this.mipmap = true
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, filter.minparam())
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, filter.maxparam())
	this.filter = filter
}

func (this *texture) SetPixels(pixels []byte) {
	if len(pixels) != this.len() {
		panic(errors.New("wrong buffer size"))
	} else {
		gl.TexSubImage2D(
			gl.TEXTURE_2D,
			0,
			0,
			0,
			int32(this.size.X),
			int32(this.size.Y),
			this.format.format(),
			gl.UNSIGNED_BYTE,
			gl.Ptr(pixels),
		)
		panicError()
	}
}

func (this *texture) Pixels(pixels []byte) []byte {
	if pixels == nil {
		pixels = make([]byte, this.len())
	} else if len(pixels) < this.len() {
		panic(errors.New("buffer too small"))
	}

	gl.GetTexImage(
		gl.TEXTURE_2D,
		0,
		this.format.format(),
		gl.UNSIGNED_BYTE,
		gl.Ptr(pixels),
	)
	panicError()

	return pixels
}

func (this *texture) len() int {
	return this.size.Y * this.size.X * this.ColorFormat().pixelSize()
}

func (this *texture) delete() {
	mainthread.CallNonBlock(func() {
		gl.DeleteTextures(1, &this.ref)
	})
}

func newTexture(size image.Point, filter Filter, format ColorFormat, pixels []byte) (*texture, error) {
	tex := &texture{
		size:   size,
		format: format,
	}
	gl.GenTextures(1, &tex.ref)
	runtime.SetFinalizer(tex, (*texture).delete)

	tex.Begin()
	defer tex.End()

	var ptr unsafe.Pointer
	if pixels != nil {
		ptr = gl.Ptr(pixels)
	} else {
		ptr = unsafe.Pointer(nil)
	}

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		format.internalFormat(),
		int32(size.X),
		int32(size.Y),
		0,
		format.format(),
		format.xtype(),
		ptr,
	)

	if err := checkError(); err != nil {
		return nil, err
	}

	tex.SetFilter(filter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	return tex, checkError()
}

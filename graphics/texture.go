package graphics

import (
	"errors"
	"image"
	"image/draw"
	"runtime"

	"github.com/askeladdk/pancake/mathx"

	gl "github.com/askeladdk/pancake/graphics/opengl"
)

const (
	TextureUnitsCount = 48 // Minimum number of Texture units according to the spec.
)

var texbinders = func(n uint32) []*binder {
	var binders []*binder
	for i := uint32(0); i < n; i++ {
		x := i // gotcha
		binders = append(binders, newBinder(func(id uint32) {
			gl.ActiveTexture(gl.Enum(gl.TEXTURE0 + x))
			gl.BindTexture(gl.TEXTURE_2D, gl.Texture(id))
		}))
	}
	return binders
}(TextureUnitsCount)

type Texture struct {
	size   image.Point
	format ColorFormat
	id     gl.Texture
}

func (tex *Texture) BeginAt(unit int) {
	texbinders[unit].bind(uint32(tex.id))
}

func (tex *Texture) EndAt(unit int) {
	texbinders[unit].unbind()
}

func (tex *Texture) Begin() {
	tex.BeginAt(0)
}

func (tex *Texture) End() {
	tex.EndAt(0)
}

func (tex *Texture) Size() image.Point {
	return tex.size
}

func (tex *Texture) ColorFormat() ColorFormat {
	return tex.format
}

func (tex *Texture) SetPixels(pixels []byte) {
	if tex.id == 0 {
		panic(errors.New("screen texture cannot be accessed"))
	} else if len(pixels) != tex.len() {
		panic(errors.New("wrong buffer size"))
	} else {
		gl.TexSubImage2D(
			gl.TEXTURE_2D,
			0,
			0,
			0,
			tex.size.X,
			tex.size.Y,
			tex.format.format(),
			gl.UNSIGNED_BYTE,
			pixels,
		)
		panicError()
	}
}

func (tex *Texture) Pixels(pixels []byte) []byte {
	if tex.id == 0 {
		panic(errors.New("screen texture cannot be accessed"))
	} else if pixels == nil {
		pixels = make([]byte, tex.len())
	} else if len(pixels) < tex.len() {
		panic(errors.New("buffer too small"))
	}

	gl.GetTexImage(
		gl.TEXTURE_2D,
		0,
		tex.format.format(),
		gl.UNSIGNED_BYTE,
		pixels,
	)
	panicError()

	return pixels
}

func (tex *Texture) Texture() *Texture {
	return tex
}

func (tex *Texture) TextureRegion() TextureRegion {
	return TextureRegion{1, 1, 0, 0}
}

func (tex *Texture) Scale() mathx.Vec2 {
	return mathx.FromPoint(tex.size)
}

func (tex *Texture) SubImage(region image.Rectangle) Image {
	return &subTexture{
		texture: tex,
		region:  NewTextureRegion(tex.Size(), region),
		size:    mathx.FromPoint(region.Size()),
	}
}

func (tex *Texture) len() int {
	return tex.size.Y * tex.size.X * tex.ColorFormat().pixelSize()
}

func (tex *Texture) delete() {
	gl.DeleteTexture(tex.id)
}

func NewTexture(size image.Point, filter Filter, format ColorFormat, pixels []byte) *Texture {
	tex := &Texture{
		id:     gl.CreateTexture(),
		size:   size,
		format: format,
	}
	runtime.SetFinalizer(tex, (*Texture).delete)

	tex.Begin()
	defer tex.End()

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		format.internalFormat(),
		size.X,
		size.Y,
		format.format(),
		gl.UNSIGNED_BYTE,
		pixels,
	)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, filter.param())
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, filter.param())
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	return tex
}

func imagePix(img image.Image) ([]uint8, ColorFormat) {
	switch im := img.(type) {
	case *image.RGBA:
		return im.Pix, ColorFormatRGBA
	case *image.Paletted:
		return im.Pix, ColorFormatIndexed
	default:
		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		return rgba.Pix, ColorFormatRGBA
	}
}

func NewTextureFromImage(img image.Image, filter Filter) *Texture {
	pix, format := imagePix(img)
	return NewTexture(img.Bounds().Size(), filter, format, pix)
}

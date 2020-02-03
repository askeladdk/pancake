package main

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/askeladdk/pancake/graphics"
	"github.com/go-gl/mathgl/mgl32"
)

type Image struct {
	Texture  *graphics.Texture
	Size     mgl32.Vec2
	UVBounds mgl32.Mat3
}

func (img *Image) SubImage(bounds image.Rectangle) *Image {
	tsize := img.Texture.Size()
	bsize := bounds.Size()
	sx := float32(bsize.X) / float32(tsize.X)
	sy := float32(bsize.Y) / float32(tsize.Y)
	tx := float32(bounds.Min.X) / float32(tsize.X)
	ty := float32(bounds.Min.Y) / float32(tsize.Y)
	return &Image{
		Texture: img.Texture,
		Size: mgl32.Vec2{
			float32(bsize.X),
			float32(bsize.Y),
		},
		UVBounds: mgl32.Mat3{
			sx, 0, 0,
			0, sy, 0,
			tx, ty, 1,
		},
	}
}

func imageToNRGBA(img image.Image) *image.NRGBA {
	switch im := img.(type) {
	case *image.NRGBA:
		return im
	default:
		nrgba := image.NewNRGBA(img.Bounds())
		draw.Draw(nrgba, nrgba.Bounds(), img, image.Point{0, 0}, draw.Src)
		return nrgba
	}
}

func NewImageFromImage(img image.Image) *Image {
	nrgba := imageToNRGBA(img)
	texture := graphics.NewTexture(nrgba.Bounds().Size(),
		graphics.FilterNearest, graphics.ColorFormatRGBA, nrgba.Pix)
	return &Image{
		Texture: texture,
		Size: mgl32.Vec2{
			float32(texture.Size().X),
			float32(texture.Size().Y),
		},
		UVBounds: mgl32.Ident3(),
	}
}

func LoadPNG(filename string) (*Image, error) {
	if f, err := os.Open(filename); err != nil {
		return nil, err
	} else if img, err := png.Decode(f); err != nil {
		return nil, err
	} else {
		return NewImageFromImage(img), nil
	}
}

package graphics

import (
	"image"

	"github.com/askeladdk/pancake/mathx"
)

type Image interface {
	Texture() *Texture
	TextureRegion() TextureRegion
	Scale() mathx.Vec2
}

type subTexture struct {
	texture *Texture
	region  TextureRegion
	size    mathx.Vec2
}

func (t *subTexture) Texture() *Texture {
	return t.texture
}

func (t *subTexture) TextureRegion() TextureRegion {
	return t.region
}

func (t *subTexture) Scale() mathx.Vec2 {
	return t.size
}

type TextureRegion struct {
	Sx, Sy, Tx, Ty float32
}

func NewTextureRegion(size image.Point, region image.Rectangle) TextureRegion {
	regionsz := region.Size()
	return TextureRegion{
		Sx: float32(regionsz.X) / float32(size.X),
		Sy: float32(regionsz.Y) / float32(size.Y),
		Tx: float32(region.Min.X) / float32(size.X),
		Ty: float32(region.Min.Y) / float32(size.Y),
	}
}

func (tr TextureRegion) Aff3() mathx.Aff3 {
	return mathx.Aff3{
		tr.Sx, 0,
		0, tr.Sy,
		tr.Tx, tr.Ty,
	}
}

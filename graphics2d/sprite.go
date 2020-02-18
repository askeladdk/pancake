package graphics2d

import (
	"image"

	"github.com/askeladdk/pancake/graphics"
	"github.com/askeladdk/pancake/mathx"
)

type Sprite struct {
	Texture *graphics.Texture
	Size    mathx.Vec2
	Region  mathx.Aff3
}

func TextureRegion(size image.Point, region image.Rectangle) (sx, sy, tx, ty float32) {
	regionsz := region.Size()
	sx = float32(regionsz.X) / float32(size.X)
	sy = float32(regionsz.Y) / float32(size.Y)
	tx = float32(region.Min.X) / float32(size.X)
	ty = float32(region.Min.Y) / float32(size.Y)
	return
}

func NewSprite(texture *graphics.Texture, region image.Rectangle) Sprite {
	sx, sy, tx, ty := TextureRegion(texture.Bounds().Size(), region)
	return Sprite{
		Texture: texture,
		Size:    mathx.FromPoint(region.Size()),
		Region: mathx.Aff3{
			sx, 0.,
			0., sy,
			tx, ty,
		},
	}
}

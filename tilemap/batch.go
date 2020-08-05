package tilemap

import (
	"image/color"

	"github.com/askeladdk/pancake/graphics"
	"github.com/askeladdk/pancake/graphics2d"
	"github.com/askeladdk/pancake/mathx"
)

func drawExtents(tileSize mathx.Vec2, viewport mathx.Rectangle) (x0, y0, x1, y1, xofs, yofs int) {
	tw, th := tileSize.Elem()
	tilew, tileh := int(tw), int(th)

	vpx0, vpy0, vpx1, vpy1 := viewport.Elem()

	xofs = -(int(vpx0) % tilew)
	yofs = -(int(vpy0) % tileh)
	xedge := int(vpx1) % tilew
	yedge := int(vpy1) % tileh
	x0 = int(vpx0) / tilew
	y0 = int(vpy0) / tileh
	x1 = int(vpx1) / tilew
	y1 = int(vpy1) / tileh

	if -xofs+xedge != 0 {
		x1 += 1
	}

	if -yofs+yedge != 0 {
		y1 += 1
	}

	return
}

// Batch draws the viewport of a TileMap. It implements the graphics2d.Batch interface.
type Batch struct {
	texture    *graphics.Texture
	regions    []graphics.TextureRegion
	modelviews []mathx.Aff3
	colors     []color.NRGBA
}

// Update the batch. Call it whenever panning the camera or when the TileMap changes.
func (b *Batch) Update(tileMap TileMap, camera *graphics2d.Camera) {
	tileSet := tileMap.TileSet()
	b.texture = tileSet.Texture()
	b.regions = b.regions[:0]
	b.modelviews = b.modelviews[:0]
	b.colors = b.colors[:0]

	tileSize := tileSet.TileSize()
	scale := mathx.ScaleAff3(tileSize)
	x0, y0, x1, y1, xofs, yofs := drawExtents(tileSize, camera.WorldViewport())
	position := mathx.Vec2{float32(xofs), float32(yofs)}.
		Add(tileSize.Mul(.5)).
		Add(camera.Viewport.Min)

	xorig := position[0]

	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if tileId := tileMap.TileAt(x, y); tileId != Absent {
				b.modelviews = append(b.modelviews, scale.Translated(position))
				b.regions = append(b.regions, tileSet.TileRegion(tileId))
				b.colors = append(b.colors, tileMap.TintColorAt(x, y))
			}
			position[0] += tileSize[0]
		}
		position[0] = xorig
		position[1] += tileSize[1]
	}
}

// Len implements graphics2d.Batch.
func (b *Batch) Len() int {
	return len(b.regions)
}

// Texture implements graphics2d.Batch.
func (b *Batch) Texture() *graphics.Texture {
	return b.texture
}

// TintColorAt implements graphics2d.Batch.
func (b *Batch) TintColorAt(i int) color.NRGBA {
	return b.colors[i]
}

// TextureRegionAt implements graphics2d.Batch.
func (b *Batch) TextureRegionAt(i int) graphics.TextureRegion {
	return b.regions[i]
}

// ModelViewAt implements graphics2d.Batch.
func (b *Batch) ModelViewAt(i int) mathx.Aff3 {
	return b.modelviews[i]
}

// OriginAt implements graphics2d.Batch.
func (b *Batch) OriginAt(i int) mathx.Vec2 {
	return mathx.Vec2{}
}

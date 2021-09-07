package tilemap

import (
	"image/color"

	"github.com/askeladdk/pancake"
	"github.com/askeladdk/pancake/mathx"
	"github.com/askeladdk/pancake/pancake2d"
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
		x1++
	}

	if -yofs+yedge != 0 {
		y1++
	}

	return
}

// Batch draws the viewport of a TileMap. It implements the graphics2d.Batch interface.
type Batch struct {
	texture    *pancake.Texture
	regions    []pancake.TextureRegion
	modelviews []mathx.Aff3
	colors     []color.Color
}

// Update the batch. Call it whenever panning the camera or when the TileMap changes.
func (b *Batch) Update(tileMap TileMap, camera *pancake2d.Camera) {
	tileSet := tileMap.TileSet()
	b.texture = tileSet.Texture()
	b.regions = b.regions[:0]
	b.modelviews = b.modelviews[:0]
	b.colors = b.colors[:0]

	tileSize := tileSet.TileSize()
	scale := mathx.ScaleAff3(tileSize)
	x0, y0, x1, y1, xofs, yofs := drawExtents(tileSize, camera.WorldViewport())
	position := mathx.Vec2{float64(xofs), float64(yofs)}.
		Add(tileSize.Mul(.5)).
		Add(camera.Viewport.Min)

	xorig := position[0]

	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if tileID := tileMap.TileAt(x, y); tileID != Absent {
				b.modelviews = append(b.modelviews, scale.Translated(position))
				b.regions = append(b.regions, tileSet.TileRegion(tileID))
				b.colors = append(b.colors, tileMap.TintColorAt(x, y))
			}
			position[0] += tileSize[0]
		}
		position[0] = xorig
		position[1] += tileSize[1]
	}
}

// Len implements pancake2d.SpriteBatch.
func (b *Batch) Len() int { return len(b.regions) }

// TextureAt implements pancake2d.SpriteBatch.
func (b *Batch) TextureAt(i int) *pancake.Texture { return b.texture }

// TintColorAt implements pancake2d.SpriteBatch.
func (b *Batch) TintColorAt(i int) color.Color { return b.colors[i] }

// TextureRegionAt implements pancake2d.SpriteBatch.
func (b *Batch) TextureRegionAt(i int) pancake.TextureRegion { return b.regions[i] }

// ModelViewAt implements pancake2d.SpriteBatch.
func (b *Batch) ModelViewAt(i int) mathx.Aff3 { return b.modelviews[i] }

// OriginAt implements pancake2d.SpriteBatch.
func (b *Batch) OriginAt(i int) mathx.Vec2 { return mathx.Vec2{} }

// ZOrderAt implements pancake2d.SpriteBatch.
func (b *Batch) ZOrderAt(i int) float64 { return 0 }

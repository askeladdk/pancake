package tilemap

import (
	"image/color"

	"github.com/askeladdk/pancake/graphics"
	"github.com/askeladdk/pancake/mathx"
)

// Batch draws the viewport of a TileMap. It implements the graphics2d.Batch interface.
type Batch struct {
	texture    *graphics.Texture
	regions    []graphics.TextureRegion
	modelviews []mathx.Aff3
	colors     []color.NRGBA
}

// Update the batch. Call it whenever panning the camera or when the TileMap changes.
func (b *Batch) Update(tileMap TileMap, camera *Camera) {
	tileset := tileMap.TileSet()
	b.texture = tileset.Texture()
	b.regions = b.regions[:0]
	b.modelviews = b.modelviews[:0]
	b.colors = b.colors[:0]
	tileMap.RangeTilesInViewport(camera.Viewport, func(cell Coordinate, modelview mathx.Aff3) {
		if tileId := tileMap.TileAt(cell); tileId != Absent {
			b.regions = append(b.regions, tileset.TileRegion(tileId))
			b.modelviews = append(b.modelviews, modelview.Translated(camera.Pos))
			b.colors = append(b.colors, tileMap.TintColorAt(cell))
		}
	})
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
	return mathx.Vec2{-.5, -.5}
}

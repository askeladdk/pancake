package tilemap

import (
	"image/color"

	"github.com/askeladdk/pancake/graphics"
	"github.com/askeladdk/pancake/mathx"
)

// Batch draws the viewport of a TileMap. It implements the graphics2d.Batch interface.
type Batch struct {
	TileMap    TileMap
	Camera     *Camera
	regions    []graphics.TextureRegion
	modelviews []mathx.Aff3
	colors     []color.NRGBA
}

// Update the batch. Call it whenever panning the camera or when the TileMap changes.
func (b *Batch) Update() {
	tileset := b.TileMap.TileSet()
	b.regions = b.regions[:0]
	b.modelviews = b.modelviews[:0]
	b.colors = b.colors[:0]
	b.TileMap.RangeTilesInViewport(b.Camera.Viewport, func(cell Coordinate, modelview mathx.Aff3) {
		if tileId := b.TileMap.TileAt(cell); tileId != Absent {
			b.regions = append(b.regions, tileset.TileRegion(tileId))
			b.modelviews = append(b.modelviews, modelview.Translated(b.Camera.Pos))
			b.colors = append(b.colors, b.TileMap.TintColorAt(cell))
		}
	})
}

// Len implements graphics2d.Batch.
func (b *Batch) Len() int {
	return len(b.regions)
}

// Texture implements graphics2d.Batch.
func (b *Batch) Texture() *graphics.Texture {
	return b.TileMap.TileSet().Texture()
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

// PivotAt implements graphics2d.Batch.
func (b *Batch) PivotAt(i int) mathx.Vec2 {
	return mathx.Vec2{.5, .5}
}

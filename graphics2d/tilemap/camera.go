package tilemap

import (
	"image"
	"image/color"

	"github.com/askeladdk/pancake/graphics"
	"github.com/askeladdk/pancake/mathx"
)

// Camera draws a region of a TileMap. It implements the graphics2d.Batch interface.
type Camera struct {
	// Pos is the top-left corner position in pixels where the TileMap will be drawn on the screen.
	Pos mathx.Vec2

	// Viewport is the area of the TileMap that is currently in view.
	// Min is the top-left corner and the difference between Min and Max is the size in pixels.
	Viewport image.Rectangle

	// TileMap is the TileMap in view.
	TileMap TileMap

	regions    []graphics.TextureRegion
	modelviews []mathx.Aff3
	colors     []color.NRGBA
}

// OnScreenArea reports the area of the screen where the viewport will be drawn.
func (b *Camera) OnScreenArea() image.Rectangle {
	x, y := b.Pos.Elem()
	return b.Viewport.Sub(b.Viewport.Min).Add(image.Pt(int(x), int(y)))
}

// Pan translates the camera viewport by delta pixels.
func (b *Camera) Pan(dp image.Point) {
	b.Viewport = mathx.ClampRectangle(b.TileMap.Bounds(), b.Viewport.Add(dp))
}

// Update the batch. Call it whenever panning the camera or when the TileMap changes.
func (b *Camera) Update() {
	tileset := b.TileMap.TileSet()
	b.regions = b.regions[:0]
	b.modelviews = b.modelviews[:0]
	b.colors = b.colors[:0]
	b.TileMap.RangeTilesInViewport(b.Viewport, func(x, y int, modelview mathx.Aff3) {
		b.regions = append(b.regions, tileset.TileRegion(b.TileMap.TileAt(x, y)))
		b.modelviews = append(b.modelviews, modelview.Translated(b.Pos))
		b.colors = append(b.colors, b.TileMap.TintColorAt(x, y))
	})
}

// Len implements graphics2d.Batch.
func (b *Camera) Len() int {
	return len(b.regions)
}

// Texture implements graphics2d.Batch.
func (b *Camera) Texture() *graphics.Texture {
	return b.TileMap.TileSet().Texture()
}

// TintColorAt implements graphics2d.Batch.
func (b *Camera) TintColorAt(i int) color.NRGBA {
	return b.colors[i]
}

// TextureRegionAt implements graphics2d.Batch.
func (b *Camera) TextureRegionAt(i int) graphics.TextureRegion {
	return b.regions[i]
}

// ModelViewAt implements graphics2d.Batch.
func (b *Camera) ModelViewAt(i int) mathx.Aff3 {
	return b.modelviews[i]
}

// PivotAt implements graphics2d.Batch.
func (b *Camera) PivotAt(i int) mathx.Vec2 {
	return mathx.Vec2{.5, .5}
}

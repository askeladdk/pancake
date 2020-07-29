package tilemap

import (
	"image"

	"github.com/askeladdk/pancake/mathx"
)

type Bounded interface {
	Bounds() image.Rectangle
}

// Camera controls the viewport of a TileMap.
type Camera struct {
	// Pos is the top-left corner position in pixels where the TileMap will be drawn on the screen.
	Pos mathx.Vec2

	// Viewport is the area of the TileMap that is currently in view.
	// Min is the top-left corner and the difference between Min and Max is the size in pixels.
	Viewport image.Rectangle

	// Bounded knows the total viewable area of the TileMap.
	Bounded Bounded
}

// OnScreenArea reports the area of the screen where the viewport will be drawn.
func (b *Camera) OnScreenArea() image.Rectangle {
	x, y := b.Pos.Elem()
	return b.Viewport.Sub(b.Viewport.Min).Add(image.Pt(int(x), int(y)))
}

// WorldToScreen converts a pixel in world coordinates to screen coordinates.
func (b *Camera) WorldToScreen(pt image.Point) mathx.Vec2 {
	return mathx.FromPoint(pt.Sub(b.Viewport.Min)).Add(b.Pos)
}

// ScreenToWorld converts a pixel in screen coordinates to world coordinates.
func (b *Camera) ScreenToWorld(v mathx.Vec2) image.Point {
	v = v.Sub(b.Pos).Add(mathx.FromPoint(b.Viewport.Min))
	return image.Pt(int(v[0]), int(v[1]))
}

// Pan translates the camera viewport by delta pixels.
func (b *Camera) Pan(dp image.Point) {
	b.Viewport = mathx.ClampRectangle(b.Bounded.Bounds(), b.Viewport.Add(dp))
}

// CenterAt centers the camera viewport at the point.
func (b *Camera) CenterAt(pt image.Point) {
	size := b.Viewport.Size()
	r := image.Rectangle{Max: size}.Add(pt.Sub(size.Div(2)))
	b.Viewport = mathx.ClampRectangle(b.Bounded.Bounds(), r)
}

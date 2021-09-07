package pancake2d

import (
	"image"

	"github.com/askeladdk/pancake/mathx"
)

// Camera controls the panning of a viewport.
type Camera struct {
	// Viewport is the area of the screen where the world will be drawn in screen coordinates.
	Viewport mathx.Rectangle

	// Bounds is the total viewable area in world coordinates.
	Bounds mathx.Rectangle

	// Pos is the current top left position in world coordinates.
	Pos mathx.Vec2

	// Target is the target top left position in world coordinates.
	Target mathx.Vec2

	// Smoothing controls the smoothing of the camera when panning.
	// It is expected to be a value between 0 and 1.
	// Higher values means faster smoothing.
	// Zero or negative means that smoothing is disabled.
	Smoothing float64
}

// Scissor reports the rectangle of the screen to scissor.
func (c *Camera) Scissor() image.Rectangle {
	x0, y0, x1, y1 := c.Viewport.Elem()
	return image.Rect(int(x0), int(y0), int(x1), int(y1))
}

// WorldToScreen converts a pixel in world coordinates to screen coordinates.
func (c *Camera) WorldToScreen(px mathx.Vec2) mathx.Vec2 {
	return px.Sub(c.Pos).Add(c.Viewport.Min)
}

// ScreenToWorld converts a pixel in screen coordinates to world coordinates.
func (c *Camera) ScreenToWorld(px mathx.Vec2) mathx.Vec2 {
	return px.Sub(c.Viewport.Min).Add(c.Pos)
}

// Frame calculates one frame of camera panning towards the target.
func (c *Camera) Frame() {
	// Snap to target when close enough.
	if c.Smoothing <= 0 || c.Target.Sub(c.Pos).LenSqr() < .5 {
		c.Pos = c.Target
	} else {
		c.Pos = c.Pos.Lerp(c.Target, mathx.Smooth(c.Smoothing))
	}
}

// Pan translates the camera viewport relative to its current position.
func (c *Camera) Pan(dt mathx.Vec2) {
	c.Target = c.Viewport.
		Sub(c.Viewport.Min).
		Add(c.Target.Add(dt)).
		Clamp(c.Bounds).
		Min
}

// CenterAt centers the camera viewport at the point.
func (c *Camera) CenterAt(px mathx.Vec2) {
	screen := c.Viewport.Sub(c.Viewport.Min)
	c.Target = screen.
		Add(px).
		Sub(screen.Size().Mul(.5)).
		Clamp(c.Bounds).
		Min
}

// WorldViewport reports the currently visible area in world coordinates.
func (c *Camera) WorldViewport() mathx.Rectangle {
	return c.Viewport.Sub(c.Viewport.Min).Add(c.Pos)
}

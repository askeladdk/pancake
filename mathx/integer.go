package mathx

import "image"

// ClampInt returns value bounded by min and max.
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

// ClampRectangle clamps a rectangle to stay within the bounds of another rectangle.
func ClampRectangle(bounds, r image.Rectangle) image.Rectangle {
	boundsSize := bounds.Size()
	size := r.Size()
	size.X = ClampInt(size.X, 0, boundsSize.X)
	size.Y = ClampInt(size.Y, 0, boundsSize.Y)
	r.Min.X = ClampInt(r.Min.X, bounds.Min.X, boundsSize.X-size.X)
	r.Min.Y = ClampInt(r.Min.Y, bounds.Min.Y, boundsSize.Y-size.Y)
	r.Max = r.Min.Add(size)
	return r
}

package mathx

import (
	"image"
	"math"
)

// Rectangle represents an axis-aligned bounding box (AABB)
// bounded by (x0, y0) -- (x1, y1).
type Rectangle struct {
	Min, Max Vec2
}

// Rect returns a canonical Rectangle in the area (x0, y0) -- (x1, y1).
func Rect(x0, y0, x1, y1 float64) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle{Vec2{x0, y0}, Vec2{x1, y1}}
}

// FromRectangle converts an image.Rectangle to a Rectangle.
func FromRectangle(p image.Rectangle) Rectangle {
	return Rectangle{
		Min: FromPoint(p.Min),
		Max: FromPoint(p.Max),
	}
}

// Elem decomposes r in its individual elements.
func (r Rectangle) Elem() (x0, y0, x1, y1 float64) {
	x0, y0 = r.Min.Elem()
	x1, y1 = r.Max.Elem()
	return
}

// Canon returns the canonical version of r
// where x0 <= x1 and y0 <= y1.
func (r Rectangle) Canon() Rectangle {
	return Rect(r.Elem())
}

// Add computes (x0 + ux, y0 + uy) -- (x1 + ux, y1 + uy).
func (r Rectangle) Add(u Vec2) Rectangle {
	r.Min = r.Min.Add(u)
	r.Max = r.Max.Add(u)
	return r
}

// Sub computes (x0 - ux, y0 - uy) -- (x1 - ux, y1 - uy).
func (r Rectangle) Sub(u Vec2) Rectangle {
	r.Min = r.Min.Sub(u)
	r.Max = r.Max.Sub(u)
	return r
}

// Dx computes x1 - x0.
func (r Rectangle) Dx() float64 {
	return r.Max[0] - r.Min[0]
}

// Dy computes y1 - y0.
func (r Rectangle) Dy() float64 {
	return r.Max[1] - r.Min[1]
}

// Size computes (x1 - x0, y1 - y0).
func (r Rectangle) Size() Vec2 {
	return r.Max.Sub(r.Min)
}

// Empty returns whether the length or width of r is 0.
func (r Rectangle) Empty() bool {
	return r.Min[0] >= r.Max[0] || r.Min[1] >= r.Max[1]
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r Rectangle) Eq(s Rectangle) bool {
	return r == s || r.Empty() && s.Empty()
}

// Expand increases the area of r to
// (x0 - v0, y0 - v1) - (x1 + v0, y1 + v1).
func (r Rectangle) Expand(v Vec2) Rectangle {
	return Rectangle{
		r.Min.Sub(v),
		r.Max.Add(v),
	}
}

// IntersectsPoint tests whether u is in r.
func (r Rectangle) IntersectsPoint(u Vec2) bool {
	return u.IntersectsRectangle(r)
}

// IntersectsCircle tests whether rectangle r and circle c intersect.
func (r Rectangle) IntersectsCircle(c Circle) bool {
	return c.IntersectsRectangle(r)
}

// IntersectsRectangle tests whether rectangles r and s have a non-empty intersection.
func (r Rectangle) IntersectsRectangle(s Rectangle) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min[0] < s.Max[0] && s.Min[0] < r.Max[0] &&
		r.Min[1] < s.Max[1] && s.Min[1] < r.Max[1]
}

// Clamp clamps rectangle r to stay within bounds without changing the size of r.
func (r Rectangle) Clamp(bounds Rectangle) Rectangle {
	boundsSize := bounds.Size()
	size := r.Size()
	size[0] = Clamp(size[0], 0, boundsSize[0])
	size[1] = Clamp(size[1], 0, boundsSize[1])
	r.Min[0] = Clamp(r.Min[0], bounds.Min[0], bounds.Max[0]-size[0])
	r.Min[1] = Clamp(r.Min[1], bounds.Min[1], bounds.Max[1]-size[1])
	r.Max = r.Min.Add(size)
	return r
}

// Circle represents a circle with radius r centered around (x, y).
type Circle struct {
	Center Vec2
	Radius float64
}

// C returns a canonical Circle with center (x, y) and radius |r|.
func C(x, y, r float64) Circle {
	return Circle{Vec2{x, y}, math.Abs(r)}
}

// Elem decomposes c in its individual elements.
func (c Circle) Elem() (x, y, r float64) {
	x, y = c.Center.Elem()
	r = c.Radius
	return
}

// Canon returns the canonical version of c with radius |r|.
func (c Circle) Canon() Circle {
	return C(c.Elem())
}

// Add computes (x + ux, y + uy).
func (c Circle) Add(u Vec2) Circle {
	c.Center = c.Center.Add(u)
	return c
}

// Add computes (x - ux, y - uy).
func (c Circle) Sub(u Vec2) Circle {
	c.Center = c.Center.Sub(u)
	return c
}

// Empty returns whether r <= 0.
func (c Circle) Empty() bool {
	return c.Radius <= 0
}

// IntersectsPoint tests whether u is in c.
func (c Circle) IntersectsPoint(u Vec2) bool {
	return u.IntersectsCircle(c)
}

// IntersectsCircle tests whether circles c0 and c1 intersect.
func (c0 Circle) IntersectsCircle(c1 Circle) bool {
	r := c0.Radius + c1.Radius
	return c1.Center.Sub(c0.Center).LenSqr() < r*r
}

// IntersectsCircle tests whether circle c intersects with rectangle r.
func (c Circle) IntersectsRectangle(r Rectangle) bool {
	closest := Vec2{
		math.Max(r.Min[0], math.Min(c.Center[0], r.Max[0])),
		math.Max(r.Min[1], math.Min(c.Center[1], r.Max[1])),
	}
	return closest.IntersectsCircle(c)
}

// Geometry intersects with geometric objects.
type Geometry interface {
	IntersectsPoint(Vec2) bool
	IntersectsCircle(Circle) bool
	IntersectsRectangle(Rectangle) bool
}

// Intersects tests whether geometries are intersecting.
func Intersects(a, b Geometry) bool {
	switch x := b.(type) {
	case *Rectangle:
		return a.IntersectsRectangle(*x)
	case *Circle:
		return a.IntersectsCircle(*x)
	case *Vec2:
		return a.IntersectsPoint(*x)
	default:
		return false
	}
}

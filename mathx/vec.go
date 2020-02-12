package mathx

import (
	"image"
	"math"

	"golang.org/x/image/math/f32"
)

// Vec2 is a 2-element vector.
type Vec2 f32.Vec2

// Vec3 is a 3-element vector.
type Vec3 f32.Vec3

// Vec4 is a 4-element vector.
type Vec4 f32.Vec4

// FromPoint converts an image.Point to a Vec2.
func FromPoint(p image.Point) Vec2 {
	return Vec2{float32(p.X), float32(p.Y)}
}

// FromHeading converts an angle in radians to a Vec2.
func FromHeading(radians float32) Vec2 {
	s, c := math.Sincos(float64(radians))
	return Vec2{float32(c), float32(s)}
}

// Vec3 converts u to a Vec3.
func (u Vec2) Vec3(z float32) Vec3 {
	return Vec3{u[0], u[1], z}
}

// Vec4 converts u to a Vec4.
func (u Vec2) Vec4(z, w float32) Vec4 {
	return Vec4{u[0], u[1], z, w}
}

// X returns u[0].
func (u Vec2) X() float32 {
	return u[0]
}

// Y returns u[1].
func (u Vec2) Y() float32 {
	return u[1]
}

// Elem decomposes u in its individual elements.
func (u Vec2) Elem() (float32, float32) {
	return u[0], u[1]
}

// Add translates u by v.
func (u Vec2) Add(v Vec2) Vec2 {
	return Vec2{u[0] + v[0], u[1] + v[1]}
}

// Sub translates u by v.
func (u Vec2) Sub(v Vec2) Vec2 {
	return Vec2{u[0] - v[0], u[1] - v[1]}
}

// Mul multiplies u by a scalar.
func (u Vec2) Mul(a float32) Vec2 {
	return Vec2{u[0] * a, u[1] * a}
}

// MulVec2 computes the element-wise multiplication of u and v.
func (u Vec2) MulVec2(v Vec2) Vec2 {
	return Vec2{u[0] * v[0], u[1] * v[1]}
}

// Neg negates u.
func (u Vec2) Neg() Vec2 {
	u[0], u[1] = -u[0], -u[1]
	return u
}

// Dot computes the dot-product of u and v.
func (u Vec2) Dot(v Vec2) float32 {
	return u[0]*v[0] + u[1]*v[1]
}

// Cross computes the cross-product of u and v.
func (u Vec2) Cross(v Vec2) float32 {
	return u[0]*v[1] - u[1]*v[0]
}

// Len computes |u|.
func (u Vec2) Len() float32 {
	return float32(math.Hypot(float64(u[0]), float64(u[1])))
}

// LenSqr computes |u|^2 which is less computationally expensive than Len.
func (u Vec2) LenSqr() float32 {
	return u.Dot(u)
}

// Normal computes the normal of u.
func (u Vec2) Normal() Vec2 {
	return Vec2{-u[1], u[0]}
}

// Unit computes u / |u|.
func (u Vec2) Unit() Vec2 {
	if u == (Vec2{}) {
		return Vec2{1, 0}
	}
	return u.Mul(1 / u.Len())
}

// IntersectsPoint tests whether u and v are (nearly) equal.
func (u Vec2) IntersectsPoint(v Vec2) bool {
	return FloatEq(u[0], v[0]) && FloatEq(u[1], v[1])
}

// IntersectsCircle tests whether u is in c.
func (u Vec2) IntersectsCircle(c Circle) bool {
	return !c.Empty() && c.Center.Sub(u).LenSqr() < c.Radius*c.Radius
}

// IntersectsRectangle tests whether u is in r.
func (u Vec2) IntersectsRectangle(r Rectangle) bool {
	return !r.Empty() && r.Min[0] <= u[0] && u[0] < r.Max[0] && r.Max[0] <= u[1] && u[1] < r.Max[1]
}

// Heading computer the angle of u in radians.
func (u Vec2) Heading() float32 {
	return float32(math.Atan2(float64(u[1]), float64(u[0])))
}

// Lerp computes the linear interpolation between u and v modulated by t.
func (u Vec2) Lerp(v Vec2, t float32) Vec2 {
	return Vec2{Lerp(u[0], v[0], t), Lerp(u[1], v[1], t)}
}

// Wrap wraps u around r.
func (u Vec2) Wrap(r Rectangle) Vec2 {
	u[0] = Wrap(u[0], r.Min[0], r.Max[0])
	u[1] = Wrap(u[1], r.Min[1], r.Max[1])
	return u
}

// Project computes the orthogonal projection of u
// onto a straight line parallel to v.
func (u Vec2) Project(v Vec2) Vec2 {
	w := v.Unit()
	return w.Mul(u.Dot(w))
}

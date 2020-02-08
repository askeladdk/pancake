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

func FromPoint(p image.Point) Vec2 {
	return Vec2{float32(p.X), float32(p.Y)}
}

func FromHeading(radians float32) Vec2 {
	s, c := math.Sincos(float64(radians))
	return Vec2{float32(c), float32(s)}
}

func (u Vec2) Vec3(z float32) Vec3 {
	return Vec3{u[0], u[1], z}
}

func (u Vec2) Vec4(z, w float32) Vec4 {
	return Vec4{u[0], u[1], z, w}
}

func (u Vec2) X() float32 {
	return u[0]
}

func (u Vec2) Y() float32 {
	return u[1]
}

func (u Vec2) Elem() (float32, float32) {
	return u[0], u[1]
}

func (u Vec2) Add(v Vec2) Vec2 {
	return Vec2{u[0] + v[0], u[1] + v[1]}
}

func (u Vec2) Sub(v Vec2) Vec2 {
	return Vec2{u[0] - v[0], u[1] - v[1]}
}

func (u Vec2) Mul(a float32) Vec2 {
	return Vec2{u[0] * a, u[1] * a}
}

func (u Vec2) MulVec2(v Vec2) Vec2 {
	return Vec2{u[0] * v[0], u[1] * v[1]}
}

func (u Vec2) Neg() Vec2 {
	u[0], u[1] = -u[0], -u[1]
	return u
}

func (u Vec2) Dot(v Vec2) float32 {
	return u[0]*v[0] + u[1]*v[1]
}

func (u Vec2) Cross(v Vec2) float32 {
	return u[0]*v[1] - u[1]*v[0]
}

func (u Vec2) Len() float32 {
	return float32(math.Hypot(float64(u[0]), float64(u[1])))
}

func (u Vec2) LenSqr() float32 {
	return u.Dot(u)
}

func (u Vec2) Normal() Vec2 {
	return Vec2{-u[1], u[0]}
}

func (u Vec2) Unit() Vec2 {
	if u == (Vec2{}) {
		return Vec2{1, 0}
	}
	return u.Mul(1 / u.Len())
}

func (u Vec2) In(r Rectangle) bool {
	return r.Min[0] <= u[0] && u[0] < r.Max[0] && r.Max[0] <= u[1] && u[1] < r.Max[1]
}

func (u Vec2) Heading() float32 {
	return float32(math.Atan2(float64(u[1]), float64(u[0])))
}

func (u Vec2) Lerp(v Vec2, t float32) Vec2 {
	return Vec2{Lerp(u[0], v[0], t), Lerp(u[1], v[1], t)}
}

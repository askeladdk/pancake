package mathx

import (
	"math"

	"golang.org/x/image/math/f32"
)

// Mat3 is a 3x3 matrix in row major order.
type Mat3 f32.Mat3

func Ident3() Mat3 {
	return Mat3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
}

// Mat4 is a 4x4 matrix in row major order.
type Mat4 f32.Mat4

func Ident4() Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Aff3 is a 3x3 affine transformation matrix in row major order,
// where the bottom row is implicitly [0 0 1].
//
//  [0] [2] [4]
//  [1] [3] [5]
//   0   0   1
type Aff3 f32.Aff3

func IdentAff3() Aff3 {
	return Aff3{
		1, 0,
		0, 1,
		0, 0,
	}
}

func TranslateAff3(u Vec2) Aff3 {
	return Aff3{
		1, 0,
		0, 1,
		u[0], u[1],
	}
}

func ScaleAff3(u Vec2) Aff3 {
	return Aff3{
		u[0], 0,
		0, u[1],
		0, 0,
	}
}

func RotateAff3(radians float32) Aff3 {
	s, c := math.Sincos(float64(radians))
	return Aff3{float32(c), float32(s), float32(-s), float32(c), 0, 0}
}

func (m Aff3) Translated(u Vec2) Aff3 {
	m[4], m[5] = m[4]+u[0], m[5]+u[1]
	return m
}

func (m Aff3) Scaled(u Vec2) Aff3 {
	m[0], m[2], m[4] = m[0]*u[0], m[2]*u[0], m[4]*u[0]
	m[1], m[3], m[5] = m[1]*u[1], m[3]*u[1], m[5]*u[1]
	return m
}

func (m Aff3) Rotated(radians float32) Aff3 {
	return m.Mul3(RotateAff3(radians))
}

func (m Aff3) Mul(v float32) Aff3 {
	return Aff3{
		v * m[0],
		v * m[1],
		v * m[2],
		v * m[3],
		v * m[4],
		v * m[5],
	}
}

func (m Aff3) Mul3(r Aff3) Aff3 {
	return Aff3{
		r[0]*m[0] + r[2]*m[1],
		r[1]*m[0] + r[3]*m[1],
		r[0]*m[2] + r[2]*m[3],
		r[1]*m[2] + r[3]*m[3],
		r[0]*m[4] + r[2]*m[5] + r[4],
		r[1]*m[4] + r[3]*m[5] + r[5],
	}
}

func (m Aff3) Mat3() Mat3 {
	return Mat3{
		m[0], m[1], 0,
		m[2], m[3], 0,
		m[4], m[5], 1,
	}
}

func (m Aff3) Mat4() Mat4 {
	return Mat4{
		m[0], m[1], 0, 0,
		m[2], m[3], 0, 0,
		m[4], m[5], 1, 0,
		0, 0, 0, 1,
	}
}

func (m Aff3) Det() float32 {
	return m[0]*m[3] - m[2]*m[1]
}

func (m Aff3) Project(u Vec2) Vec2 {
	return Vec2{m[0]*u[0] + m[2]*u[1] + m[4], m[1]*u[0] + m[3]*u[1] + m[5]}
}

func (m Aff3) Unproject(u Vec2) Vec2 {
	det := m.Det()
	return Vec2{
		m[3]*(u[0]-m[4]) - m[2]*(u[1]-m[5]),
		-m[1]*(u[0]-m[4]) + m[0]*(u[1]-m[5]),
	}.Mul(1 / det)
}

func (m Aff3) Inv() Aff3 {
	det := m.Det()

	if FloatEq(det, 0) {
		return Aff3{}
	}

	return Aff3{
		m[3], -m[1],
		-m[2], m[0],
		m[2]*m[5] - m[3]*m[4],
		m[1]*m[4] - m[0]*m[5],
	}.Mul(1 / det)
}

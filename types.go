package pancake

import (
	"github.com/go-gl/mathgl/mgl64"
)

var quad = []float32{
	// x, y, u, v
	-1, -1, 0, 0,
	-1, +1, 0, 1,
	+1, -1, 1, 0,
	+1, +1, 1, 1,
}

type Rect [2]mgl64.Vec2

func (this Rect) Min() mgl64.Vec2 {
	return this[0]
}

func (this Rect) Max() mgl64.Vec2 {
	return this[1]
}

func (this Rect) Center() mgl64.Vec2 {
	return this[1].Add(this[0]).Mul(0.5)
}

func (this Rect) Size() mgl64.Vec2 {
	return this[1].Sub(this[0])
}

func (this Rect) Elem() (x0, y0, x1, y1 float64) {
	x0, y0 = this[0].Elem()
	x1, y1 = this[1].Elem()
	return
}

func R(x0, y0, x1, y1 float64) Rect {
	return Rect{mgl64.Vec2{x0, y0}, mgl64.Vec2{x1, y1}}
}

package mathx

type Rectangle struct {
	Min, Max Vec2
}

func Rect(x0, y0, x1, y1 float32) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle{Vec2{x0, y0}, Vec2{x1, y1}}
}

func (r Rectangle) Elems() (float32, float32, float32, float32) {
	return r.Min[0], r.Min[1], r.Max[0], r.Max[1]
}

func (r Rectangle) Dx() float32 {
	return r.Max[0] - r.Min[0]
}

func (r Rectangle) Dy() float32 {
	return r.Max[1] - r.Min[1]
}

func (r Rectangle) Size() Vec2 {
	return r.Max.Sub(r.Min)
}

func (r Rectangle) Empty() bool {
	return r.Min[0] >= r.Max[0] || r.Min[1] >= r.Max[1]
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r Rectangle) Eq(s Rectangle) bool {
	return r == s || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rectangle) Overlaps(s Rectangle) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min[0] < s.Max[0] && s.Min[0] < r.Max[0] &&
		r.Min[1] < s.Max[1] && s.Min[1] < r.Max[1]
}

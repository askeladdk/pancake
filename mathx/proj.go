package mathx

func Ortho(left, right, bottom, top, near, far float32) Mat4 {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)

	return Mat4{
		2 / rml, 0, 0, 0,
		0, 2 / tmb, 0, 0,
		0, 0, -2 / fmn, 0,
		-(right + left) / rml, -(top + bottom) / tmb, -(far + near) / fmn, 1,
	}
}

func Ortho2D(left, right, bottom, top float32) Mat4 {
	return Ortho(left, right, bottom, top, -1, 1)
}

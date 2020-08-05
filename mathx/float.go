package mathx

import "math"

const (
	epsilon = 1e-10
	// mask     = 0x7F
	// shift    = 32 - 8 - 1
	signMask = 1 << 31
	// uvnan    = 0x7F800001
	// uvinf    = 0x7F000000
	// uvneginf = 0xFF000000

	// Tau is 2pi, a full turn.
	Tau = 2 * math.Pi
)

// // IsNaN reports whether f is an IEEE 754 ``not-a-number'' value.
// func IsNaN(f float32) (is bool) {
// 	// return f != f
// 	x := math.Float32bits(f)
// 	return uint32(x>>shift)&mask == mask && x != uvinf && x != uvneginf
// }

// // IsInf reports whether f is an infinity, according to sign.
// // If sign > 0, IsInf reports whether f is positive infinity.
// // If sign < 0, IsInf reports whether f is negative infinity.
// // If sign == 0, IsInf reports whether f is either infinity.
// func IsInf(f float32, sign int) bool {
// 	// return sign >= 0 && f > math.MaxFloat32 || sign <= 0 && f < -math.MaxFloat32
// 	x := math.Float32bits(f)
// 	return sign >= 0 && x == uvinf || sign <= 0 && x == uvneginf
// }

// func NaN() float32 { return math.Float32frombits(uvnan) }

// func Inf(sign int) float32 {
// 	if sign >= 0 {
// 		return math.Float32frombits(uvinf)
// 	} else {
// 		return math.Float32frombits(uvneginf)
// 	}
// }

// func Signbit(x float32) bool {
// 	return math.Float32bits(x)&signMask != 0
// }

// Min returns the smaller value given x and y.
func Min(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

// Min returns the greater value given x and y.
func Max(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}

// Abs returns the absolute value of x.
func Abs(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) &^ signMask)
}

// Clamp returns x bounded by min and max.
func Clamp(x, min, max float32) float32 {
	if x < min {
		return min
	} else if x >= max {
		return max
	}
	return x
}

// Wrap returns x wrapped around min and max.
func Wrap(x, min, max float32) float32 {
	if x < min {
		return max
	} else if x >= max {
		return min
	}
	return x
}

// FloatEq compares a and b for (near) equality,
// accounting for floating point inaccuracies.
//
// https://floating-point-gui.de/errors/comparison/
func FloatEq(a, b float32) bool {
	if a == b {
		return true
	}

	diff := Abs(a - b)
	if a*b == 0 || diff < math.SmallestNonzeroFloat32 {
		return diff < epsilon*epsilon
	}

	return diff/(Abs(a)+Abs(b)) < epsilon
}

// Lerp computes the linear interpolation of a and b modulated by t.
func Lerp(a, b, t float32) float32 {
	return a + t*(b-a)
}

// Smooth computes a smooth interpolation of t.
func Smooth(t float32) float32 {
	return t * t * (3 - 2*t)
}

// Mod returns the floating-point remainder of x/y.
func Mod(x, y float32) float32 {
	return float32(math.Mod(float64(x), float64(y)))
}

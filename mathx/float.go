package mathx

import "math"

// Tau is 2pi, a full turn.
const Tau = 2 * math.Pi

// Clamp returns x bounded by min and max.
func Clamp(x, min, max float64) float64 {
	if x < min {
		return min
	} else if x >= max {
		return max
	}
	return x
}

// Wrap returns x wrapped around min and max.
func Wrap(x, min, max float64) float64 {
	if x < min {
		return max
	} else if x >= max {
		return min
	}
	return x
}

// Equal compares a and b for (near) equality,
// accounting for floating point inaccuracies.
//
// https://floating-point-gui.de/errors/comparison/
func Equal(a, b float64) bool {
	const epsilon = 1e-10

	if a == b {
		return true
	}

	diff := math.Abs(a - b)
	if a*b == 0 || diff < math.SmallestNonzeroFloat64 {
		return diff < epsilon*epsilon
	}

	return diff/(math.Abs(a)+math.Abs(b)) < epsilon
}

// Lerp computes the linear interpolation of a and b modulated by t.
func Lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// InvLerp inverts a lerp.
//  InvLerp(a, b, Lerp(a, b, t)) = t
func InvLerp(a, b, c float64) float64 {
	return 1 - (b-c)/(b-a)
}

// Smooth computes a smooth interpolation of t.
func Smooth(t float64) float64 {
	return t * t * (3 - 2*t)
}

// Saturate computes the proportion of x between a and b as a value between 0 and 1.
//  Saturate(3, 9,  1) = 0
//  Saturate(3, 9,  6) = 0.5
//  Saturate(3, 9, 10) = 1
func Saturate(a, b, x float64) float64 {
	return Clamp((x-a)/(b-a), 0, 1)
}

// Step returns 0 if x < edge, else 1.
func Step(edge, x float64) float64 {
	if x < edge {
		return 0
	}
	return 1
}

package geom

// Bernstein basis polynomials for deg=3

func b30(x float64) float64 {
	c := 1 - x
	return c * c * c
}

func b31(x float64) float64 {
	c := 1 - x
	return 3 * x * c * c
}

// equivalent to b30 plus b31
func b30pb31(x float64) float64 {
	// 2x^3 - 3x^2 + 1
	return 2*x*x*x - (3 * x * x) + 1
}

func b32(x float64) float64 {
	return 3 * x * x * (1 - x)
}

func b33(x float64) float64 {
	return x * x * x
}

// equivalent to b32 plus b33
func b32pb33(x float64) float64 {
	// x^2 (3 - 2x)
	return x * x * (3 - 2*x)
}

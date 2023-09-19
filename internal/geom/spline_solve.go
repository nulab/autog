package geom

import "math"

func solve3(coeff []float64) []float64 {
	a, b, c, d := coeff[3], coeff[2], coeff[1], coeff[0]
	if aeq0(a) {
		return solve2(coeff)
	}
	b_over_3a := b / (3 * a)
	c_over_a := c / a
	d_over_a := d / a

	p := b_over_3a * b_over_3a
	q := 2*b_over_3a*p - b_over_3a*c_over_a + d_over_a
	p = c_over_a/3 - p
	disc := q*q + 4*p*p*p

	var roots []float64
	if disc < 0 {
		r := 0.5 * math.Sqrt(-disc+q*q)
		theta := math.Atan2(math.Sqrt(-disc), -q)
		temp := 2 * math.Cbrt(r)
		roots = []float64{
			temp * math.Cos(theta/3),
			temp * math.Cos((theta+2*math.Pi)/3),
			temp * math.Cos((theta-2*math.Pi)/3),
		}
	} else {
		alpha := 0.5 * (math.Sqrt(disc) - q)
		beta := -q - alpha
		c := math.Cbrt(alpha) + math.Cbrt(beta)
		if disc > 0 {
			roots = []float64{c}
		} else {
			roots = []float64{
				c,
				-0.5 * c,
				-0.5 * c,
			}
		}
	}
	for i := range roots {
		roots[i] -= b_over_3a
	}
	return roots
}

func solve2(coeff []float64) []float64 {
	// cubic coefficient is zero, solve as quadratic
	a, b, c := coeff[2], coeff[1], coeff[0]
	if aeq0(a) {
		// quadratic coefficient is zero, solve as linear
		return solve1(coeff)
	}
	b_over_2a := b / (2 * a)
	c_over_a := c / a

	// discriminant
	disc := b_over_2a*b_over_2a - c_over_a
	switch {
	case disc < 0:
		// no solutions in R
		return []float64{}
	case disc > 0:
		// original C code here optimizes by calculating the second root from the first
		// as derived from the quadratic equation itself basically to calculate the square root only ocne
		u := -b_over_2a + math.Sqrt(disc)
		return []float64{
			u,
			-2*b_over_2a - u,
		}
	default:
		// disc == 0, one solution
		return []float64{-b_over_2a}
	}
}

func solve1(coeff []float64) []float64 {
	a, b := coeff[1], coeff[0]
	if aeq0(a) {
		if aeq0(b) {
			// degenerate, infinite solutions
			return nil
		} else {
			// only constant term non-zero, no solution
			return []float64{}
		}
	}
	// degenerate as cx + d = 0, i.e. a line, one solution
	return []float64{-b / a}
}

func aeq0(x float64) bool {
	return x < epsilon3 && x > -epsilon3
}

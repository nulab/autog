package geom

import "math"

// MakeSpline computes control points of a cubic bezier from a to b so that it gently curves in the
// direction of the slope.
func MakeSpline(a, b P) ctrlp {
	// slide factor
	const k = 0.2
	// if the points are vertically aligned, let the cubic degenerate into a straight segment
	if a.X == b.X {
		return ctrlp{a, a, b, b}
	}
	// compute vector
	v := P{b.X - a.X, b.Y - a.Y}

	// rotation angle
	// rotate clockwise when the x of the vector is positive, counterclockwise when is negative
	// x is inverted to account for rotation in the SVG plane
	theta1 := sign(-v.X) * math.Pi * (1.0 / 10.0)
	theta2 := sign(-v.X) * math.Pi * (9.0 / 10.0)

	// rotate the vector
	r := rotate(v, theta1)
	q := rotate(v, theta2)
	// slide the control points along the rotated vector
	p1 := P{a.X + k*r.X, a.Y + k*r.Y}
	p2 := P{b.X + k*q.X, b.Y + k*q.Y}

	return ctrlp{
		p0: a,
		p1: p1,
		p2: p2,
		p3: b,
	}
}

func sign(x float64) float64 {
	if x < 0 {
		return -1
	}
	return 1
}

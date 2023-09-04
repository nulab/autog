package geom

import (
	"fmt"
	"math"
)

// P represents a point on the plane
type P struct {
	X, Y float64
}

func (p P) String() string {
	return fmt.Sprintf(`<circle r="4" cx="%.02f" cy="%.02f" fill="black"/>`, p.X, p.Y)
}

// could also be represented as a [2]float64 but there's essentially no difference except for readability: p.X, p.Y vs p[0], p[1]
// otherwise both are comparable, take up 16 bytes and have meaningful zero values

// adds p2 to p1 and returns a new point
func addp(p1, p2 P) P {
	return P{
		p1.X + p2.X,
		p1.Y + p2.Y,
	}
}

// subtracts p2 from p1 and returns a new point
func subp(p1, p2 P) P {
	return P{
		p1.X - p2.X,
		p1.Y - p2.Y,
	}
}

// multiplies p by a scalar and returns a new point
func scalep(p P, c float64) P {
	return P{
		p.X * c,
		p.Y * c,
	}
}

// computes the dot product between p1 and p2
func dotp(p1, p2 P) float64 {
	return p1.X*p2.X + p1.Y*p2.Y
}

// computes the distance between p1 and p2
func distp(p, q P) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// computes the square distance between p1 and p2
func sqdistp(p, q P) float64 {
	return (q.X-p.X)*(q.X-p.X) + (q.Y-p.Y)*(q.Y-p.Y)
}

// normalizes the vector represented by this point (sets its length to 1)
func norm(p P) P {
	d := p.X*p.X + p.Y*p.Y
	if d > epsilon2 {
		d = math.Sqrt(d)
		return P{
			X: p.X / d,
			Y: p.Y / d,
		}
	}
	return p
}

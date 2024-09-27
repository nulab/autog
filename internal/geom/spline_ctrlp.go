package geom

import "fmt"

// cubic bezier control points
type ctrlp struct {
	p0 P // start control point
	p1 P // first movable control point
	p2 P // second movable control point
	p3 P // end control point
}

func (bz ctrlp) String() string {
	return fmt.Sprintf(`<path d="M %.02f %.02f C %.02f %.02f, %.02f %.02f, %.02f %.02f" stroke="black"/>`,
		bz.p0.X, bz.p0.Y, bz.p1.X, bz.p1.Y, bz.p2.X, bz.p2.Y, bz.p3.X, bz.p3.Y)
}

func (bz ctrlp) Float64Slice() [][2]float64 {
	return [][2]float64{
		{bz.p0.X, bz.p0.Y},
		{bz.p1.X, bz.p1.Y},
		{bz.p2.X, bz.p2.Y},
		{bz.p3.X, bz.p3.Y},
	}
}

// Returns a new set of control points with:
//   - P0 and P3 unchanged
//   - P1 and P2 moved along their tangents by thirds of the slide factor
func (bz ctrlp) adjust(a float64) ctrlp {
	return ctrlp{
		p0: bz.p0,
		p1: addp(bz.p0, scalep(bz.p1, a/3.0)),
		p2: subp(bz.p3, scalep(bz.p2, a/3.0)),
		p3: bz.p3,
	}
}

func (bz ctrlp) dist() float64 {
	return distp(bz.p0, bz.p1) + distp(bz.p1, bz.p2) + distp(bz.p2, bz.p3)
}

func (bz ctrlp) xcoeff() []float64 {
	return bz.coeff(bz.p0.X, bz.p1.X, bz.p2.X, bz.p3.X)
}

func (bz ctrlp) ycoeff() []float64 {
	return bz.coeff(bz.p0.Y, bz.p1.Y, bz.p2.Y, bz.p3.Y)
}

func (bz ctrlp) scoeff(slope float64) []float64 {
	v0 := bz.p0.Y - slope*bz.p0.X
	v1 := bz.p1.Y - slope*bz.p1.X
	v2 := bz.p2.Y - slope*bz.p2.X
	v3 := bz.p3.Y - slope*bz.p3.X
	return bz.coeff(v0, v1, v2, v3)
}

// coefficients of the polynomial form of the cubic bezier
func (bz ctrlp) coeff(v0, v1, v2, v3 float64) []float64 {
	return []float64{
		v0,
		3 * (v1 - v0),
		3*v0 + 3*v2 - 6*v1,
		v3 + 3*v1 - (v0 + 3*v2),
	}
}

// Returns the point on the parametric curve that corresponds to the given t
func (bz ctrlp) curvep(t float64) P {
	return P{
		X: b30(t)*bz.p0.X + b31(t)*bz.p1.X + b32(t)*bz.p2.X + b33(t)*bz.p3.X,
		Y: b30(t)*bz.p0.Y + b31(t)*bz.p1.Y + b32(t)*bz.p2.Y + b33(t)*bz.p3.Y,
	}
}

// Returns the index of the path point farthest from the curve
func (bz ctrlp) maxerr(path []P, t []float64) int {
	if len(path) != len(t) {
		panic("spline: path and parameters have unequal length")
	}
	maxd := -1.0
	maxi := -1

	for i := 1; i < len(path)-1; i++ {
		d := distp(bz.curvep(t[i]), path[i])
		if d > maxd {
			maxd = d
			maxi = i
		}
	}
	return maxi
}

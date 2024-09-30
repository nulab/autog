package geom

import (
	"math"
)

const (
	numthresh = 1e-6 // threshold for numerical problems
)

// FitSpline computes the control points of a cubic bezier curve that fits a polygonal path.
// The input path points must be ordered from start to end.
func FitSpline(path []P, tanv1, tanv2 P, barriers []Segment) []ctrlp {
	// This function implements P. Schneider's "An algorithm for automatically fitting digitized curves".
	// The underlying logic is summarized as follows:
	// given a set of points that form a polygonal path, the first and last points are set equal to
	// the first and last control points P0 and P3 of the cubic bezier; the remaining controls P1 and P2 must be placed along
	// the tangents at P0 and P3 in order to preserve continuity.
	// In case of a single bezier segment, there is nothing to "continue" from, however tangents are used to
	// enforce exit angles. In case of a piece-wise curve, this ensures continuity at the joints.
	// The algorithm defines the P1 and P2 control points in terms of their collinear counterparts P0 and P3,
	// specifically P1 = alpha1 * tangent1 + P0 and P2 = alpha2 * tangent2 + P3
	// where alpha1 and alpha2 are the two scale factors by which to move P1 and P2 along the respective tangents.

	tanv1 = norm(tanv1)
	tanv2 = norm(tanv2)

	// curve parametrization
	t := chordLength(path)
	// compute A terms, which depend on the start and end tangent vectors and the curve parameters
	a := aterms(tanv1, tanv2, t)

	// find the scale factors of the other two control points
	alpha1, alpha2 := alphas(t, a, path)

	v1 := scalep(tanv1, alpha1)
	v2 := scalep(tanv2, alpha2)

	// from now on the code is essentially a port of Graphviz Dot C implementation
	// reference paper: "Implementing a General-Purpose Edge Router" by Dobkin, Gansner et al., section 3.2
	// https://dpd.cs.princeton.edu/Papers/DGKN97.pdf

	// initial control points of cubic bezier
	bz := ctrlp{
		p0: path[0],
		p1: v1,
		p2: v2,
		p3: path[len(path)-1],
	}

	bz, ok := tryfit(bz, path, barriers)
	if ok {
		return []ctrlp{bz}
	}

	// reset the movable control points to a shorter distance from the start and end control points
	bz = bz.adjust(1)
	// find the index of the point of maximum distance from the path, i.e. the point of maximum error of the curve fitting
	k := bz.maxerr(path, t)

	// compute the direction of the curve at the split point
	tansplit1 := norm(subp(path[k], path[k-1]))
	tansplit2 := norm(subp(path[k+1], path[k]))
	// average out the directions at the split point
	tanw := norm(addp(tansplit1, tansplit2))

	// recursively fit the spline to the upper half and lower half of the path
	// using tanw ensures curve continuity at the split point by making sure that
	// upper P2, upper P3, lower P0, lower P1 are collinear, in particular upper P3 = lower P0
	upperps := FitSpline(path[:k+1], tanv1, tanw, barriers)
	lowerps := FitSpline(path[k:], tanw, tanv2, barriers)

	return append(upperps, lowerps...)
}

func chordLength(path []P) []float64 {
	t := make([]float64, len(path))
	t[0] = 0
	// compute cumulative distance
	for i := 1; i < len(path); i++ {
		t[i] = t[i-1] + distp(path[i], path[i-1])
	}
	// normalize
	for i := 1; i < len(path); i++ {
		t[i] /= t[len(path)-1]
	}
	return t
}

// returns the A(i,1) and A(i,2) terms used in Schneider's method
func aterms(tanv1, tanv2 P, t []float64) [][2]P {
	a := make([][2]P, len(t))
	for i := 0; i < len(t); i++ {
		a[i][0] = scalep(tanv1, b31(t[i]))
		a[i][1] = scalep(tanv2, b32(t[i]))
	}
	return a
}

// solve the partial derivatives in alpha1 and alpha2 of the objective function
// defined as the sum of square distances between points of the polygonal path and the corresponding points on the parametric curve
func alphas(t []float64, a [][2]P, path []P) (alpha1, alpha2 float64) {
	// Schneider's paper stores them in the 2x2 matrix called C
	// since the dot product is commutative and c(1,2) = c(2,1) here I use three variables
	a1sq := 0.0 // square of A(i,1) = c(1,1)
	a2sq := 0.0 // square of A(i,2) = c(2,2)
	a1a2 := 0.0 // A(i,1) dot A(i,2)

	// init the X terms
	x1 := 0.0
	x2 := 0.0

	// summations resulting in the alpha coefficients and X terms
	for i := range t {
		a1sq += dotp(a[i][0], a[i][0])
		a1a2 += dotp(a[i][0], a[i][1])
		a2sq += dotp(a[i][1], a[i][1])

		d := path[i]
		// V0 (B30 + B31) + V3 (B32 + B33)
		Q := addp(
			scalep(path[0], b30pb31(t[i])),
			scalep(path[len(path)-1], b32pb33(t[i])),
		)
		diffdQ := subp(d, Q)

		x1 += dotp(diffdQ, a[i][0])
		x2 += dotp(diffdQ, a[i][1])
	}

	dc1x := a1sq*x2 - x1*a1a2   // determinant of [c11 c12][x1 x2]
	dxc2 := x1*a2sq - a1a2*x2   // determinant of [x1 x2][c12 c22]
	dc := a1sq*a2sq - a1a2*a1a2 // determinant of the C matrix

	// here the implementation of Graphviz Dot checks for numerical problems
	// due to very small float values; I just replicate the fallback
	if math.Abs(dc) >= numthresh {
		// proceed with Schneider's final step
		alpha1 = dxc2 / dc
		alpha2 = dc1x / dc
	}
	// in case of numerical problems, or if the numerators went to zero (happens when the polygonal path is just one segment)
	if math.Abs(dc) < numthresh || alpha1 <= 0 || alpha2 <= 0 {
		// set alphas to a third of the distance between first and last points
		k := distp(path[0], path[len(path)-1]) / 3.0
		alpha1 = k
		alpha2 = k
	}
	return
}

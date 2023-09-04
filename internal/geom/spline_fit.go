package geom

const (
	epsilon1    = 1e-3
	epsilon2    = 1e-6
	epsilon3    = 1e-7
	slideinit   = 3.0
	slidethresh = 0.01  // threshold above which we keep reducing the distance
	slidemin    = 0.005 // threshold below which we stop iterating
)

// Tests whether the curve identified by the given control points fits within a set of edges.
// If it doesn't fit, the control points P1 and P2 are moved closer to P0 and P3 until P0=P1 and P3=P2.
// If the routine finds a set of fitting control points, it returns those control points and true.
// Otherwise, it returns the original control points and false.
func tryfit(bz0 ctrlp, path []P, barriers []Segment) (ctrlp, bool) {
	slide := slideinit
	pathd := 0.0
	for i := 1; i < len(path); i++ {
		pathd += distp(path[i], path[i-1])
	}
	pathn := len(path)
	first := true

	for {
		bz1 := bz0.adjust(slide)

		// at the first iteration, check if the total traveled distance between the control points of the curve
		// is somehow shorter than the shortest path it fits to
		if first {
			if bz1.dist() < pathd-epsilon1 {
				return bz0, false
			}
			first = false
		}

		if curveContained(bz1, barriers) {
			return bz1, true
		}

		if slide < slidemin {
			// when the factor reaches zero, the cubic becomes a straight segment
			// if the input path has only two points, it's equivalent to the curve: we consider this a success
			if pathn == 2 {
				return bz1, true
			}
			return bz0, false
		}
		// move P1 and P2 closer to P0 and P3 by halving the distance
		// or drop to zero after a certain threshold
		if slide > slidethresh {
			slide /= 2
		} else {
			slide = 0
		}
	}
}

func curveContained(bz ctrlp, barriers []Segment) bool {
	for _, b := range barriers {
		roots := curveIntersects(bz, b)
		if roots == nil {
			continue
		}
		for _, r := range roots {
			if r < epsilon2 || r > 1-epsilon2 {
				continue
			}
			rp := bz.curvep(r)

			if sqdistp(rp, b.A) < epsilon1 || sqdistp(rp, b.B) < epsilon1 {
				continue
			}
			return false
		}
	}
	return true
}

func curveIntersects(bz ctrlp, seg Segment) (roots []float64) {
	xc0 := seg.A.X
	xc1 := seg.B.X - seg.A.X
	yc0 := seg.A.Y
	yc1 := seg.B.Y - seg.A.Y

	appendr01 := func(r float64) {
		if r >= 0 && r <= 1 {
			roots = append(roots, r)
		}
	}

	if xc1 == 0 {
		if yc1 == 0 {
			// here the segment degenerates into a point
			curvexc := bz.xcoeff()
			curvexc[0] -= xc0
			xroots := solve3(curvexc)

			curveyc := bz.ycoeff()
			curveyc[0] -= yc0
			yroots := solve3(curveyc)

			if xroots == nil {
				if yroots == nil {
					return nil
				} else {
					for _, yr := range yroots {
						appendr01(yr)
					}
				}
			} else if yroots == nil {
				for _, xr := range xroots {
					appendr01(xr)
				}
			} else {
				for _, xr := range xroots {
					for _, yr := range yroots {
						if xr == yr {
							appendr01(xr)
						}
					}
				}
			}
			return roots
		}
		// xc1 == 0, yc1 != 0 then the segment is vertical
		curvexc := bz.xcoeff()
		curvexc[0] -= xc0
		xroots := solve3(curvexc)
		if xroots == nil {
			return nil
		}
		for _, tv := range xroots {
			if tv >= 0 && tv <= 1 {
				curveyc := bz.ycoeff()
				sv := curveyc[0] + tv*(curveyc[1]+tv*(curveyc[2]+tv*curveyc[3]))
				sv = (sv - yc0) / yc1
				if 0 <= sv && sv <= 1 {
					appendr01(tv)
				}
			}
		}
		return roots
	} else {
		slope := yc1 / xc1
		curvesc := bz.scoeff(slope)
		curvesc[0] += slope*xc0 - yc0

		xroots := solve3(curvesc)
		if xroots == nil {
			return nil
		}
		for _, tv := range xroots {
			if tv >= 0 && tv <= 1 {
				curvexc := bz.xcoeff()
				sv := curvexc[0] + tv*(curvexc[1]+tv*(curvexc[2]+tv*curvexc[3]))
				sv = (sv - xc0) / xc1
				if 0 <= sv && sv <= 1 {
					appendr01(tv)
				}
			}
		}
		return roots
	}
}

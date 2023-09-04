package geom

import "slices"

// Triangulate returns a list of triangles resulting from the triangulation of the polygon obtained by
// merging the rectangles (the rectangles aren't actually merged).
// Since the polygon's internal angles are all fractions of π, the triangulation is special-cased to run in O(n).
// Note that the polygon is y-monotone but isn't also strictly y-monotone, therefore the left and right vertex chains
// can't be considered sorted.
// For each output triangle, the segment (A,B) is a diagonal of the polygon.
func Triangulate(rects []Rect) []Tri {
	id := 0
	nextid := func() int {
		id++
		return id
	}

	if len(rects) == 1 {
		return []Tri{
			{nextid(), rects[0].BR, rects[0].TL, P{rects[0].BR.X, rects[0].TL.Y}},
			{nextid(), rects[0].BR, rects[0].TL, P{rects[0].TL.X, rects[0].BR.Y}},
		}
	}

	ts := make([]Tri, 0, len(rects)*4)
	for i := 0; i < len(rects); i++ {
		r1 := rects[i]

		var a, b P
		if i < len(rects)-1 {
			r2 := rects[i+1]
			a, b = left2right(r1.BR, P{r2.BR.X, r2.TL.Y})
		} else {
			a = r1.BR
		}

		hasMergePoint := false
		if i > 0 {
			r0 := rects[i-1]
			if r1.TL.X < r0.TL.X {
				// merge point on left chain
				s := P{X: r0.TL.X, Y: r1.TL.Y}
				c := leftmost(r0.BR, P{r1.BR.X, r1.TL.Y})
				ts = append(ts, Tri{nextid(), a, s, c})
				ts = append(ts, Tri{nextid(), a, r1.TL, s})
				hasMergePoint = true
			}
			if r0.BR.X < r1.BR.X {
				// merge point on right chain
				s := P{X: r0.BR.X, Y: r1.TL.Y}
				ts = append(ts, Tri{nextid(), a, s, P{r1.BR.X, r1.TL.Y}})
				ts = append(ts, Tri{nextid(), a, r1.TL, s})
				hasMergePoint = true
			}
			if i == len(rects)-1 {
				ts = append(ts, Tri{nextid(), r1.BR, r1.TL, P{r1.TL.X, r1.BR.Y}})
				if !hasMergePoint {
					ts = append(ts, Tri{nextid(), r1.TL, r1.BR, P{r1.BR.X, r1.TL.Y}})
				}
			}

		}
		if i < len(rects)-1 {
			r2 := rects[i+1]
			if r1.BR.X > r2.BR.X {
				// split point on right chain
				s := P{X: r1.BR.X, Y: r1.TL.Y}
				ts = append(ts, Tri{nextid(), a, s, b})
			}
			// horizontal diagonal from right to left chain
			ts = append(ts, Tri{nextid(), a, rightmost(P{X: r1.TL.X, Y: r1.BR.Y}, r2.TL), r1.TL})
			if !hasMergePoint {
				ts = append(ts, Tri{nextid(), r1.TL, a, P{r1.BR.X, r1.TL.Y}})
			}

			if r1.TL.X < r2.TL.X {
				// split point on left chain
				ts = append(ts, Tri{nextid(), r2.TL, r1.TL, P{r1.TL.X, r1.BR.Y}})
			}
		}
	}

	return slices.Clip(ts)
}

func rightmost(p1, p2 P) P {
	if p1.X < p2.X {
		return p2
	}
	return p1
}

func leftmost(p1, p2 P) P {
	if p1.X < p2.X {
		return p1
	}
	return p2
}

func left2right(p1, p2 P) (P, P) {
	if p1.X < p2.X {
		return p1, p2
	}
	return p2, p1
}

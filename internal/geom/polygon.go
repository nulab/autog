package geom

// Polygon represents a polygon on the plane
type Polygon struct {
	Points []P // counterclockwise list of points
	r      int // index at which the right chain starts
}

func (p Polygon) Sides() []Segment {
	barriers := []Segment{}
	for i := 1; i < len(p.Points); i++ {
		barriers = append(barriers, Segment{p.Points[i-1], p.Points[i]})
	}
	barriers = append(barriers, Segment{p.Points[len(p.Points)-1], p.Points[0]})
	return barriers
}

func MergeRects(rects []Rect) Polygon {
	np := len(rects) * 2 // number of polygon vertices

	lps := make([]P, 0, np)
	rps := make([]P, 0, np)

	var prev Rect
	for i, r := range rects {
		if i == 0 {
			lps = append(lps, r.TL)
			rps = append(rps, P{r.BR.X, r.TL.Y})
			prev = r
			continue
		}

		if prev.TL.X != r.TL.X {
			lps = append(lps, P{prev.TL.X, r.TL.Y}, r.TL)
		} else {
			lps = append(lps, r.TL)
		}

		if prev.BR.X != r.BR.X {
			// rps will be iterated backwards
			rps = append(rps, prev.BR, P{r.BR.X, prev.BR.Y})
		} else {
			rps = append(rps, prev.BR)
		}

		prev = r
	}
	lps = append(lps, P{prev.TL.X, prev.BR.Y})
	rps = append(rps, prev.BR)

	points := make([]P, 0, len(lps)+len(rps))
	i := 0
	for i < len(lps) {
		points = append(points, lps[i])
		i++
	}
	r := i
	j := len(rps) - 1
	for j >= 0 {
		points = append(points, rps[j])
		i++
		j--
	}
	return Polygon{points, r}
}

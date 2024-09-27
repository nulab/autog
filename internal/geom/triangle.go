package geom

import "fmt"

// Tri represents a triangle
type Tri struct {
	ID      int
	A, B, C P
}

func (t Tri) SVG() string {
	return fmt.Sprintf(`<path d="M %.2f,%.2f %.2f,%.2f %.2f,%.2f Z" stroke="blue" fill="none"/>`, t.A.X, t.A.Y, t.B.X, t.B.Y, t.C.X, t.C.Y)
}

func (t Tri) Barycenter() P {
	bx := (t.A.X + t.B.X + t.C.X) / 3
	by := (t.A.Y + t.B.Y + t.C.Y) / 3
	return P{bx, by}
}

func (t Tri) Contains(p P) bool {
	s := 0
	e := []P{t.A, t.B, t.C}
	for i := range 3 {
		or := orientation(e[i%3], e[(i+1)%3], p)
		if or == cln {
			q, r := e[i%3], e[(i+1)%3]
			return p.X >= min(q.X, r.X) && p.X <= max(q.X, r.X) &&
				p.Y >= min(q.Y, r.Y) && p.Y <= max(q.Y, r.Y)
		}

		if or != cw {
			s++
		}
	}
	return s == 3 || s == 0
}

func (t Tri) OrderedSide(i int) Segment {
	e := []P{t.A, t.B, t.C}
	a, b := e[i%3], e[(i+1)%3]
	if a.X < b.X {
		return Segment{a, b}
	}
	if b.X < a.X {
		return Segment{b, a}
	}
	if a.Y < b.Y {
		return Segment{a, b}
	}
	return Segment{b, a}
}

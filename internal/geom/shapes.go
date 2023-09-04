package geom

import "fmt"

type Segment struct {
	A, B P
}

func (seg Segment) String() string {
	return fmt.Sprintf(`<path d="M %.2f,%.2f %.2f,%.2f" stroke="blue" />`, seg.A.X, seg.A.Y, seg.B.X, seg.B.Y)
}

func (seg Segment) Other(v P) P {
	if seg.A == v {
		return seg.B
	}
	return seg.A
}

// Tri represents a triangle
type Tri struct {
	ID      int
	A, B, C P
}

func (t Tri) String() string {
	return fmt.Sprintf(`<path d="M %.2f,%.2f %.2f,%.2f %.2f,%.2f Z" stroke="blue" fill="lightgrey"/>`, t.A.X, t.A.Y, t.B.X, t.B.Y, t.C.X, t.C.Y)
	// return fmt.Sprintf(`<path d="M %.2f,%.2f %.2f,%.2f" stroke="blue"/>`, t.A.X, t.A.Y, t.B.X, t.B.Y)
}

func (t Tri) Barycenter() P {
	bx := (t.A.X + t.B.X + t.C.X) / 3
	by := (t.A.Y + t.B.Y + t.C.Y) / 3
	return P{bx, by}
}

func (t Tri) Contains(p P) bool {
	s := 0
	e := []P{t.A, t.B, t.C}
	for i := 0; i < 3; i++ {
		if orientation(e[i%3], e[(i+1)%3], p) != cw {
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

// Rect represents a rectangle
type Rect struct {
	TL P // top-left vertex of the rectangle
	BR P // bottom-right vertex of the rectangle
}

func (r Rect) String() string {
	return fmt.Sprintf("{geom.P{%.02f,%.02f},geom.P{%.02f,%.02f}}", r.TL.X, r.TL.Y, r.BR.X, r.BR.Y)
}

func (r Rect) Width() float64 {
	return r.BR.X - r.TL.X
}

func (r Rect) Height() float64 {
	return r.BR.Y - r.TL.Y
}

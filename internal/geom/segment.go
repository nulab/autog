package geom

import "fmt"

type Segment struct {
	A, B P
}

func (seg Segment) SVG() string {
	return fmt.Sprintf(`<path d="M %.2f,%.2f %.2f,%.2f" stroke="blue" />`, seg.A.X, seg.A.Y, seg.B.X, seg.B.Y)
}

func (seg Segment) Other(v P) P {
	if seg.A == v {
		return seg.B
	}
	return seg.A
}

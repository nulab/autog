package geom

import "fmt"

// Rect represents a rectangle
type Rect struct {
	TL P // top-left vertex of the rectangle
	BR P // bottom-right vertex of the rectangle
}

func (r Rect) String() string {
	return fmt.Sprintf("{geom.P{%.02f,%.02f},geom.P{%.02f,%.02f}}", r.TL.X, r.TL.Y, r.BR.X, r.BR.Y)
}

func (r Rect) SVG() string {
	width := r.BR.X - r.TL.X
	height := r.BR.Y - r.TL.Y
	return fmt.Sprintf(`<rect class="rect" x="%f" y="%f" width="%f" height="%f" style="fill: lightgrey; stroke: black;" />`, r.TL.X, r.TL.Y, width, height)
}

func (r Rect) Width() float64 {
	return r.BR.X - r.TL.X
}

func (r Rect) Height() float64 {
	return r.BR.Y - r.TL.Y
}

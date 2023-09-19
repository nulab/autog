package geom

const (
	// counter-clockwise
	ccw = iota - 1
	// collinear
	cln
	// clockwise
	cw
)

// calculates the determinant of the vector product of the three points, and determines the orientation.
// NOTE: autog works with SVG-like coordinates so the inequalities are reversed
func orientation(a, b, c P) int {
	d := (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
	if d < 0 {
		return ccw
	}
	if d > 0 {
		return cw
	}
	return cln
}

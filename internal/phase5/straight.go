package phase5

import "github.com/nulab/autog/internal/graph"

func execStraightRouting(routes []routableEdge) {
	for _, r := range routes {
		if r.IsFlat() {
			r.Points = flatStraight(r.ns[0], r.ns[len(r.ns)-1])
		} else {
			r.Points = straight(r.ns[0], r.ns[len(r.ns)-1])
		}
	}
}

func flatStraight(from, to *graph.Node) [][2]float64 {
	// middle of right side
	x1 := from.X + from.W
	y1 := from.Y + from.H/2
	// middle of left side
	x2 := to.X
	y2 := to.Y + to.H/2
	// return points
	return [][2]float64{{x1, y1}, {x2, y2}}
}

func straight(from, to *graph.Node) [][2]float64 {
	// middle of lower side
	x1 := from.X + from.W/2
	y1 := from.Y + from.H
	// middle of upper side
	x2 := to.X + to.W/2
	y2 := to.Y
	// return points
	return [][2]float64{{x1, y1}, {x2, y2}}
}

func flatStartPoint(e *graph.Edge) [2]float64 {
	var x, y float64
	if e.From.LayerPos < e.To.LayerPos {

	} else {
		// middle of left side
		x = e.From.X
		y = e.From.Y + e.From.H/2
	}
	return [2]float64{x, y}
}

func flatEndPoint(e *graph.Edge) [2]float64 {
	var x, y float64
	if e.From.LayerPos < e.To.LayerPos {

	} else {
		// middle of right side
		x = e.To.X + e.To.W
		y = e.To.Y + e.To.H/2
	}
	return [2]float64{x, y}
}

func startPoint(e *graph.Edge) [2]float64 {
	var x, y float64
	if e.From.Layer < e.To.Layer {

	} else {
		// middle of upper side
		x = e.From.X + e.From.W/2
		y = e.From.Y
	}
	return [2]float64{x, y}
}

func endPoint(e *graph.Edge) [2]float64 {
	var x, y float64
	if e.From.Layer < e.To.Layer {
		// middle of upper side
		x = e.To.X + e.To.W/2
		y = e.To.Y
	} else {
		// middle of lower side
		x = e.To.X + e.To.W/2
		y = e.To.Y + e.To.H
	}
	return [2]float64{x, y}
}

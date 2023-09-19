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

func straight(from, to *graph.Node) [][2]float64 {
	return [][2]float64{startPoint(from), endPoint(to)}
}

// middle of lower side
func startPoint(n *graph.Node) [2]float64 {
	x := n.X + n.W/2
	y := n.Y + n.H
	return [2]float64{x, y}
}

// middle of upper side
func endPoint(to *graph.Node) [2]float64 {
	x := to.X + to.W/2
	y := to.Y
	return [2]float64{x, y}
}

package phase5

import (
	"github.com/nulab/autog/internal/graph"
)

func execPolylineRouting(g *graph.DGraph, routes []routableEdge) {
	for _, r := range routes {
		if r.IsFlat() {
			flatPolyline(r, g.Layers[r.From.Layer].H)
			continue
		}
		if len(r.ns) == 2 {
			r.Points = straight(r.ns[0], r.ns[len(r.ns)-1])
			continue
		}
		r.Points = append(r.Points, startPoint(r.ns[0]))
		for _, n := range r.ns[1 : len(r.ns)-1] {
			r.Points = append(r.Points, nonTerminalPoint(n, g.Layers[n.Layer].H))
		}
		r.Points = append(r.Points, endPoint(r.ns[len(r.ns)-1]))
	}
}

func nonTerminalPoint(n *graph.Node, layerHeight float64) [2]float64 {
	if !n.IsVirtual {
		panic("routing: bend point on non-virtual node")
	}
	var x, y float64
	x = n.X + n.W/2
	y = n.Y + layerHeight/2
	return [2]float64{x, y}
}

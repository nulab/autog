package phase5

import (
	"github.com/nulab/autog/internal/graph"
	"github.com/nulab/autog/internal/num"
)

func execPolylineRouting(g *graph.DGraph, routes []routableEdge) {
	for _, r := range routes {
		if r.IsFlat() {
			if num.Abs(r.From.LayerPos-r.To.LayerPos) > 1 {
				r.Points = flatNonConsecutive(r.Edge, g.Layers[r.From.Layer].H)
			} else {
				r.Points = flatStraight(r.ns[0], r.ns[len(r.ns)-1])
			}
		} else {
			if len(r.ns) == 2 {
				r.Points = straight(r.ns[0], r.ns[len(r.ns)-1])
			} else {
				r.Points = append(r.Points, startPoint(r.ns[0]))
				for _, n := range r.ns[1 : len(r.ns)-1] {
					r.Points = append(r.Points, nonTerminalPoint(n, g.Layers[n.Layer].H))
				}
				r.Points = append(r.Points, endPoint(r.ns[len(r.ns)-1]))
			}

		}
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

func flatNonConsecutive(e *graph.Edge, layerH float64) [][2]float64 {
	anchorXOffset := 20.0
	dist := num.Abs(e.From.LayerPos - e.To.LayerPos)
	startx := e.From.X + e.From.W
	starty := e.From.Y + e.From.H/2
	endx := e.To.X
	endy := e.To.Y + e.To.H/2
	top := min(starty-layerH/2, endy-layerH/2) - (10 + float64(dist)*5)

	points := [][2]float64{
		{startx, starty},
		{startx + anchorXOffset, starty},
		{startx + anchorXOffset, top},
		{endx - anchorXOffset, top},
		{endx - anchorXOffset, endy},
		{endx, endy},
	}
	return points
}

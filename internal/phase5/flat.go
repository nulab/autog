package phase5

import (
	"github.com/nulab/autog/internal/graph"
	"github.com/nulab/autog/internal/num"
)

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

func flatPolyline(r routableEdge, layerH float64) {
	if num.Abs(r.From.LayerPos-r.To.LayerPos) > 1 {
		r.Points = flatNonConsecutive(r.Edge, layerH)
	} else {
		r.Points = flatStraight(r.ns[0], r.ns[len(r.ns)-1])
	}
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

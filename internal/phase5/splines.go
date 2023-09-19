package phase5

import (
	"github.com/nulab/autog/internal/geom"
	"github.com/nulab/autog/internal/graph"
)

// todo: work in progress
func execSplines(g *graph.DGraph, params graph.Params) {
	for _, e := range g.Edges {
		rects := []geom.Rect{
			// todo: build rects
		}

		poly := geom.MergeRects(rects)

		start := geom.P{e.From.X + e.From.W/2, e.From.Y + e.From.H}
		end := geom.P{e.To.X + e.To.W/2, e.To.Y}

		path := geom.Shortest(start, end, rects)
		ctrls := geom.FitSpline(path, geom.P{}, geom.P{}, poly.Sides())
		_ = ctrls
	}
}

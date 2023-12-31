package phase5

import (
	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	imonitor.PrefixFor(alg)

	// side effects: this call merges long edges, basically undoes phase 3's breakLongEdges.
	// virtual nodes which the edges go through are collected into route structs
	// after this call, graph traversals that follow directed edges won't see virtual nodes anymore
	routableEdges := mergeLongEdges(g)

	switch alg {
	case NoRouting:
		return
	case Straight:
		execStraightRouting(routableEdges)
	case Polyline:
		execPolylineRouting(g, routableEdges)
	case Ortho:
		execOrthoRouting(g, routableEdges, params)
	default:
		panic("routing: unknown alg value")
	}

	// post-processing, restore all reversed edges
	for _, e := range g.Edges {
		if e.IsReversed {
			// reverse back
			e.Reverse()
		}
	}
}

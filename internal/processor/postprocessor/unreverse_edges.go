package postprocessor

import "github.com/nulab/autog/internal/graph"

// UnreverseEdges restores edge direction that was reversed during the cycle-breaking phase
func UnreverseEdges(g *graph.DGraph) {
	for _, e := range g.Edges {
		if e.IsReversed {
			// reverse back
			e.Reverse()
		}
	}
}

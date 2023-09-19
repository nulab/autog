package phase3

import (
	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

// Process runs this ordering algorithm on the input graph. The graph nodes must be layered.
func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	imonitor.PrefixFor(alg)
	switch alg {
	case NoOrdering:
		return
	case WMedian:
		execWeightedMedian(g, params)
	default:
		panic("ordering: unknown alg value")
	}
}

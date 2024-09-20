package phase3

import (
	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

// Process runs this ordering algorithm on the input graph. The graph nodes must be layered.
func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	imonitor.PrefixFor(alg)

	if len(g.Nodes) == 1 {
		// a single node defaults to position zero in layer zero
		imonitor.Log(imonitor.KeySkip, "not enough nodes")
		return
	}

	switch alg {
	case NoOrdering:
		return
	case WMedian:
		execWeightedMedian(g, params)
	default:
		panic("ordering: unknown alg value")
	}
}

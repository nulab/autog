package phase2

import (
	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

// Process runs this layering algorithm on the input graph. The graph must be acyclic.
func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	imonitor.PrefixFor(alg)

	if len(g.Nodes) == 1 {
		// a single node defaults to layer zero
		goto initLayers
	}

	switch alg {
	case LongestPath:
		execLongestPath(g)
	case NetworkSimplex:
		execNetworkSimplex(g, params)
	default:
		panic("layering: unknown alg value")
	}

initLayers:
	size := 0
	for _, n := range g.Nodes {
		size = max(size, n.Layer)
	}
	size += 1 // the highest layer index must fit in the slice too

	ls := make([]*graph.Layer, size)

	for _, n := range g.Nodes {
		l := ls[n.Layer]
		if l == nil {
			l = &graph.Layer{Index: n.Layer}
		}
		l.Nodes = append(l.Nodes, n)
		ls[n.Layer] = l
	}
	g.Layers = ls

	// fill layers
	for i := range size {
		l := g.Layers[i]
		if l == nil {
			g.Layers[i] = &graph.Layer{Index: i}
		}
	}
}

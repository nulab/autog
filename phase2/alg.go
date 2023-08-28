package phase2

import (
	"github.com/nulab/autog/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

type Alg uint8

func (alg Alg) Phase() int {
	return 2
}

const (
	// LongestPath computes a partition of the graph in layers by traversing nodes in topological order.
	// It may result in more flat edges and comparatively more virtual nodes, therefore more long edges too, but runs in O(N).
	// Suitable for graphs with few "flow" paths.
	LongestPath Alg = iota

	// NetworkSimplex computes a partition of the graph in layers by minimizing total edge length.
	// It results in few virtual nodes and usually no flat edges, but runs in Θ(VE). Worst case seems to be O(V^2*E)
	NetworkSimplex
)

func (alg Alg) String() (s string) {
	switch alg {
	case LongestPath:
		s = "longest path"
	case NetworkSimplex:
		s = "network simplex"
	default:
		s = "<invalid>"
	}
	return
}

func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	imonitor.PrefixFor(alg)
	switch alg {
	case LongestPath:
		execLongestPath(g)
	case NetworkSimplex:
		execNetworkSimplex(g, params)
	default:
		panic("layering: unknown alg value")
	}

	m := map[int]*graph.Layer{}
	for _, n := range g.Nodes {
		layer := m[n.Layer]
		if layer == nil {
			layer = &graph.Layer{Index: n.Layer}
		}
		layer.Nodes = append(layer.Nodes, n)
		m[n.Layer] = layer
	}
	g.Layers = m
	fillLayers(g)
}

func fillLayers(g *graph.DGraph) {
	highest := 0
	for i := range g.Layers {
		highest = max(highest, i)
	}
	for i := 0; i < highest; i++ {
		_, ok := g.Layers[i]
		if !ok {
			g.Layers[i] = &graph.Layer{Index: i}
		}
	}
}

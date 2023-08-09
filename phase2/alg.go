package phase2

import (
	"github.com/nulab/autog/graph"
)

type Alg uint8

const (
	// LongestPath computes a partition of the graph in layers by traversing nodes in topological order.
	// It may result in more flat edges and comparatively more virtual nodes, therefore more long edges too, but runs in O(N).
	// Suitable for graphs with few "flow" paths.
	LongestPath Alg = iota

	// NetworkSimplex computes a partition of the graph in layers by minimizing total edge length.
	// It results in few virtual nodes and usually no flat edges, but runs in Θ(VE). Worst case seems to be O(V^2*E)
	NetworkSimplex
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	switch alg {
	case LongestPath:
		execLongestPath(g)
		defer fillLayers(g)
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
}

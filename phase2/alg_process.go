package phase2

import (
	"github.com/nulab/autog/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

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

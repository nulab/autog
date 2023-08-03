package phase2

import (
	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/monitor"
)

type Alg uint8

const (
	NetworkSimplex Alg = iota
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph, _ *monitor.Monitor) {
	switch alg {
	case NetworkSimplex:
		execNetworkSimplex(g)
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

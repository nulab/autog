package layering

import (
	"github.com/nulab/autog/graph"
)

type Alg uint8

const (
	NetworkSimplex Alg = iota
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph) {
	switch alg {
	case NetworkSimplex:
		execNetworkSimplex(g)
	default:
		panic("layering: unknown alg value")
	}

	// todo: might abstract this into a method, behind a user option
	// 	unflattening makes the diagram strictly hierarchical but increases num of long edges
loop:
	for _, e := range g.Edges {
		if !e.SelfLoops() && e.From.Layer == e.To.Layer {
			e.To.Layer++
			goto loop
		}
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

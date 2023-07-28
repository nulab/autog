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

	// todo: it might be interesting to test whether forcing no flat edges empirically improves the layout
	// 	however right now this code (apparently) breaks the final cacoo/shape JSON output.
	// loop:
	// 	loop := true
	// 	for loop {
	// 		loop = false
	// 		for _, e := range g.Edges {
	// 			if e.From.Layer == e.To.Layer {
	//
	// 				g.Layers[e.To.Layer].RemoveNode(e.To)
	// 				e.To.Layer++
	// 				nextl := g.Layers[e.To.Layer]
	// 				if nextl == nil {
	// 					nextl = &graph.Layer{Index: e.To.Layer}
	// 				}
	// 				nextl.Nodes = append(nextl.Nodes, e.To)
	// 				g.Layers[e.To.Layer] = nextl
	//
	// 				goto loop
	// 			}
	// 		}
	// 	}
}

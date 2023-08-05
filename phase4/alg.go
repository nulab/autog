package phase4

import (
	"github.com/nulab/autog/graph"
)

type Alg uint8

const (
	// NoPositioning does nothing. Nodes won't be assigned any coordinates.
	NoPositioning Alg = iota

	// VerticalAlign aligns nodes in each layer around the center of the diagram.
	// It's a simple and fast to implement algorithm for quick prototyping.
	VerticalAlign

	// BrandesKoepf aligns nodes based on blocks and classes. Runs in O(V+E).
	BrandesKoepf

	// NetworkSimplex sets X coordinates by constructing an auxiliary graph and solving it with the network simplex method.
	// Layers in the auxiliary graph are X coordinates in the main graph.
	// Time-intensive for large graphs, above a few dozen nodes.
	NetworkSimplex

	// SinkColoring is a variant of BrandesKöpf that aligns nodes based on vertical blocks starting from the bottom.
	SinkColoring
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	switch alg {
	case NoPositioning:
		return
	case VerticalAlign:
		execVerticalAlign(g, params)
	case BrandesKoepf:
		execBrandesKoepf(g, params)
	case NetworkSimplex:
		execNetworkSimplex(g, params)
	case SinkColoring:
		execSinkColoring(g, params)
	default:
		panic("positioning: unknown alg value")
	}
	assignYCoords(g, params.LayerSpacing)
}

func assignYCoords(g *graph.DGraph, layerSpacing float64) {
	y := 0.0
	for i := 0; i < len(g.Layers); i++ {
		for _, n := range g.Layers[i].Nodes {
			n.Y = y
		}
		y += g.Layers[i].H + layerSpacing
	}
}

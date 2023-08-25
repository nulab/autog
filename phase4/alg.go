package phase4

import (
	"github.com/nulab/autog/graph"
)

type Alg uint8

const (
	// NoPositioning does nothing. Nodes won't be assigned any coordinates.
	NoPositioning Alg = iota

	// VerticalAlign aligns nodes in each layer vertically around the center of the diagram.
	// Works best for tree-like graphs with no back-edges.
	VerticalAlign

	MedianAlign

	// BrandesKoepf aligns nodes based on blocks and classes in O(V+E).
	// It results in a compact drawing but with less long straight edges.
	BrandesKoepf

	// NetworkSimplex sets X coordinates by constructing an auxiliary graph and solving it with the network simplex method.
	// Layers in the auxiliary graph are X coordinates in the main graph. Might be time-intensive for graphs above a few dozen nodes.
	NetworkSimplex

	// SinkColoring is a variant of BrandesKöpf that aligns nodes based on vertical blocks starting from the bottom.
	// It results in a larger drawing but with more long vertical edge paths. Runs in O(2kn) with 1 <= k <= maxshifts.
	SinkColoring

	// PackRight aligns nodes to the right.
	PackRight
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	switch alg {
	case NoPositioning:
		// no-op, but assign Y coordinates
	case VerticalAlign:
		execVerticalAlign(g, params)
	case MedianAlign:
		execMedianAlign(g, params)
	case BrandesKoepf:
		execBrandesKoepf(g, params)
	case NetworkSimplex:
		execNetworkSimplex(g, params)
	case SinkColoring:
		execSinkColoring(g, params)
	case PackRight:
		execPackRight(g, params)
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

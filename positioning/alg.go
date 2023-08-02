package positioning

import (
	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/monitor"
)

type Alg uint8

const (
	// NoPositioning does nothing. Nodes won't be assigned any coordinates.
	NoPositioning Alg = iota

	// VerticalAlign aligns nodes in each layer around the center of the diagram.
	// It's a simple and fast to implement algorithm for quick prototyping.
	VerticalAlign

	// BrandesKoepfExtended aligns nodes based on blocks and classes. Runs in O(V+E).
	BrandesKoepfExtended

	// NetworkSimplex sets X coordinates by constructing an auxiliary graph and solving it with the network simplex method.
	// Layers in the auxiliary graph are X coordinates in the main graph. Time-intensive for large graphs. Gansner et al.
	// mention graph size above "a few dozen" nodes.
	NetworkSimplex
	_endAlg
)

const (
	defaultLayerSpacing = 150
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph, _ *monitor.Monitor) {
	switch alg {
	case NoPositioning:
		return
	case VerticalAlign:
		execVerticalAlign(g)
	case BrandesKoepfExtended:
		execBrandesKoepf(g)
	case NetworkSimplex:
		execNetworkSimplex(g)
	default:
		panic("positioning: unknown alg value")
	}
	assignYCoords(g)
}

func assignYCoords(g *graph.DGraph) {
	y := 0.0
	for i := 0; i < len(g.Layers); i++ {
		for _, n := range g.Layers[i].Nodes {
			n.Y = y
		}
		y += g.Layers[i].H + defaultLayerSpacing
	}
}

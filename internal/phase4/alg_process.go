package phase4

import (
	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

// Process runs this positioning algorithm on the input graph. The graph nodes must be layered and ordered.
func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	imonitor.PrefixFor(alg)
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

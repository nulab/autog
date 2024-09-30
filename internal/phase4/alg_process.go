package phase4

import (
	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

// Process runs this positioning algorithm on the input graph. The graph nodes must be layered and ordered.
func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	imonitor.PrefixFor(alg)

	if len(g.Nodes) == 1 {
		// the node defaults to position (0,0) and the layer is as large as the node itself
		g.Layers[0].W = g.Nodes[0].W
		g.Layers[0].H = g.Nodes[0].H
		return
	}

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
	for _, l := range g.Layers {
		for _, n := range l.Nodes {
			n.Y = y
		}
		y += l.H + layerSpacing
	}
}

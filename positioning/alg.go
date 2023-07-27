package positioning

import (
	"github.com/nulab/autog/graph"
)

type Alg uint8

const (
	VerticalAlign Alg = iota
	BrandesKoepfExtended
	NetworkSimplex
	_endAlg
)

const (
	defaultLayerSpacing = 150
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph) {
	g.LayersX = make(map[int]*graph.Layer)
	for k, nodes := range g.Layers {
		g.LayersX[k] = &graph.Layer{
			Nodes: nodes,
			Index: k,
		}
	}

	switch alg {
	case VerticalAlign:
		execVerticalAlign(g)
	case BrandesKoepfExtended:
		execBrandesKoepf(g)
	case NetworkSimplex:
		panic("not implemented")
	default:
		panic("positioning: unknown alg value")
	}
	assignYCoords(g)
}

func assignYCoords(g *graph.DGraph) {
	y := 0.0
	for i := 0; i < len(g.LayersX); i++ {
		for _, n := range g.LayersX[i].Nodes {
			n.Y = y
		}
		y += g.LayersX[i].H + defaultLayerSpacing
	}
}

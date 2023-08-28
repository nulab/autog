package phase5

import (
	"github.com/nulab/autog/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

const (
	edgeTypeNoneVirtual = iota
	edgeTypeOneVirtual
	edgeTypeBothVirtual
)

// todo: improve code reuse of routing algos

func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	imonitor.PrefixFor(alg)
	switch alg {
	case NoRouting:
		return
	case Straight:
		execStraightRouting(g)
	case PieceWise:
		execPieceWiseRouting(g)
	case Ortho:
		execOrthoRouting(g, params)
	default:
		panic("routing: unknown alg value")
	}
}

package phase3

import (
	"github.com/nulab/autog/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	imonitor.PrefixFor(alg)
	switch alg {
	case NoOrdering:
		return
	case GraphvizDot:
		execGraphvizDot(g, params)
	default:
		panic("ordering: unknown alg value")
	}
}

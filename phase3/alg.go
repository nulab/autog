package phase3

import (
	"github.com/nulab/autog/graph"
)

type Alg uint8

const (
	// NoOrdering does nothing. Nodes won't be reordered in their layers to minimize edge crossings.
	NoOrdering Alg = iota

	// GraphvizDot implements the mincross heuristic used in dot. It attempts to minimize bilayer edge crossings
	// by sweeping up and down the layers and applying ordering nodes based on their weighted medians.
	GraphvizDot
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	switch alg {
	case NoOrdering:
		return
	case GraphvizDot:
		execGraphvizDot(g, params)
	default:
		panic("ordering: unknown alg value")
	}
}

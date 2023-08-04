package phase1

import (
	"github.com/nulab/autog/graph"
)

type Alg uint8

const (
	Greedy Alg = iota // todo: document that this is non-deterministic
	DepthFirst
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph, params graph.Params) {

	switch alg {
	case Greedy:
		execGreedy(g)
	case DepthFirst:
		execDepthFirst(g)
	default:
		panic("cyclebreaking: unknown alg value")
	}
}

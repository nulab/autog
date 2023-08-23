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
	removeTwoNodeCycles(g)
	if !hasCycles(g) {
		return
	}
	switch alg {
	case Greedy:
		execGreedy(g)
	case DepthFirst:
		execDepthFirst(g)
	default:
		panic("cyclebreaking: unknown alg value")
	}
}

func removeTwoNodeCycles(g *graph.DGraph) {
	type pair [2]*graph.Node

	seen := map[pair]bool{}
	rev := graph.EdgeSet{}

	for _, e := range g.Edges {
		a, b := e.From, e.To
		if seen[pair{a, b}] || seen[pair{b, a}] {
			rev[e] = true
		} else {
			seen[pair{a, b}] = true
		}
	}
	for e := range rev {
		e.Reverse()
	}
}

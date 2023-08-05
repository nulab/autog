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
	hide := graph.EdgeSet{}

	for _, e := range g.Edges {
		a, b := e.From, e.To
		if seen[pair{a, b}] || seen[pair{b, a}] {
			hide[e] = true
		} else {
			seen[pair{a, b}] = true
		}
	}
	for e := range hide {
		e.From.Out.Remove(e)
		e.To.In.Remove(e)
		g.Edges.Remove(e)
		g.HiddenEdges.Add(e)
	}
}

package phase1

import (
	"github.com/nulab/autog/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

// Process runs this cycle breaking algorithm on the input graph. The graph must be connected.
func (alg Alg) Process(g *graph.DGraph, _ graph.Params) {
	imonitor.PrefixFor(alg)
	// preprocessing
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

	if hasCycles(g) {
		panic("cyclebreaking: graph is still cyclic")
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

package cyclebreaking

import (
	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/monitor"
)

type Alg uint8

const (
	Greedy Alg = iota
	DepthFirst
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph, _ *monitor.Monitor) {

	switch alg {
	case Greedy:
		execGreedy(g)
	case DepthFirst:
		execDepthFirst(g)
	default:
		panic("cyclebreaking: unknown alg value")
	}
}

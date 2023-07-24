package cyclebreaking

import "github.com/nulab/autog/graph"

type Alg uint8

const (
	Greedy Alg = iota
	DepthFirst
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph) {
	switch alg {
	case Greedy:
		execGreedy(g)
	case DepthFirst:
		execDepthFirst(g)
	default:
		panic("cyclebreaking: unknown enum value")
	}
}

package cyclebreaking

import "github.com/vibridi/autog/graph"

type Alg uint8

const (
	GREEDY Alg = iota
	DEPTH_FIRST
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph) {
	switch alg {
	case GREEDY:
		execGreedy(g)
	case DEPTH_FIRST:
		execDepthFirst(g)
	default:
		panic("cyclebreaking: unknown enum value")
	}
}

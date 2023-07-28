package ordering

import "github.com/nulab/autog/graph"

type Alg uint8

const (
	GraphvizDot Alg = iota
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph) {
	switch alg {
	case GraphvizDot:
		execGraphvizDot(g)
	default:
		panic("ordering: unknown alg value")
	}
}

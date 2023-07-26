package ordering

import "github.com/nulab/autog/graph"

type Alg uint8

const (
	GansnerNorth Alg = iota
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph) {
	switch alg {
	case GansnerNorth:
		execGansnerNorth(g)
	default:
		panic("ordering: unknown alg value")
	}
}

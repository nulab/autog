package layering

import "github.com/vibridi/autog/graph"

type Alg uint8

const (
	NETWORK_SIMPLEX Alg = iota
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph) {
	switch alg {
	case NETWORK_SIMPLEX:
		execNetworkSimplex(g)
	default:
		panic("layering: unknown enum value")
	}
}

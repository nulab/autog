package layering

import "github.com/nulab/autog/graph"

type Alg uint8

const (
	NetworkSimplex Alg = iota
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph) {
	switch alg {
	case NetworkSimplex:
		execNetworkSimplex(g)
	default:
		panic("layering: unknown alg value")
	}

	m := map[int][]*graph.Node{}
	for _, n := range g.Nodes {
		m[n.Layer] = append(m[n.Layer], n)
	}
	g.Layers = m
}

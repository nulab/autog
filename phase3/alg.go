package phase3

import (
	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/monitor"
)

type Alg uint8

const (
	GraphvizDot Alg = iota
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph, m *monitor.Monitor) {
	switch alg {
	case GraphvizDot:
		execGraphvizDot(g, m)
	default:
		panic("ordering: unknown alg value")
	}
}

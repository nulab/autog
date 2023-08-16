package phase5

import (
	"github.com/nulab/autog/graph"
)

type Alg uint8

const (
	// NoRouting does not compute edge points.
	NoRouting Alg = iota

	// Straight computes the start and end point of each edge. With only two points, edges can be drawn as straight lines.
	// Unsuitable for graphs with many long edges or flat edges between non-consecutive nodes as edges may overlap nodes.
	Straight

	// PieceWise computes the start, end and bend points of each edge. Bend point coordinates are where virtual nodes would be.
	// Edges can be drawn as polylines or with curved elbows if bend points are considered bezier control points.
	PieceWise
	_endAlg
)

const (
	edgeTypeNoneVirtual = iota
	edgeTypeOneVirtual
	edgeTypeBothVirtual
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph, _ graph.Params) {
	// restore any
	for _, e := range g.HiddenEdges {
		g.Edges.Add(e)
	}
	g.HiddenEdges = nil

	switch alg {
	case NoRouting:
		return
	case Straight:
		execStraightRouting(g)
	case PieceWise:
		execPieceWiseRouting(g)
	default:
		panic("routing: unknown alg value")
	}
}

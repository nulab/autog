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

	// Ortho draws edges as piecewise orthogonal segments, i.e. all edges bend at 90 degrees.
	// Dense graphs look tidier, but it's harder to understand where edges start and finish.
	// Suitable when there's few sets of edges with the same target node.
	Ortho
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

// todo: improve code reuse of routing algos

func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	switch alg {
	case NoRouting:
		return
	case Straight:
		execStraightRouting(g)
	case PieceWise:
		execPieceWiseRouting(g)
	case Ortho:
		execOrthoRouting(g, params)
	default:
		panic("routing: unknown alg value")
	}
}

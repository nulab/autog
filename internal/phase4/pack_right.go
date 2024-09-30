package phase4

import (
	"slices"

	"github.com/nulab/autog/internal/graph"
)

func execPackRight(g *graph.DGraph, params graph.Params) {
	leftBound := 0.0
	for _, l := range g.Layers {
		x := 0.0
		for _, n := range slices.Backward(l.Nodes) {
			x -= n.W + params.NodeSpacing
			n.X = x
		}
		leftBound = min(leftBound, x)
	}

	for _, l := range g.Layers {
		for _, n := range l.Nodes {
			n.X -= leftBound
			l.H = max(l.H, n.H)
		}
	}
}

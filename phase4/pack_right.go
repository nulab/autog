package phase4

import "github.com/nulab/autog/graph"

func execPackRight(g *graph.DGraph, params graph.Params) {
	leftBound := 0.0
	for _, l := range g.Layers {
		x := 0.0
		iter := nodesIterator(l.Nodes, left)
		for n := iter(); n != nil; n = iter() {
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

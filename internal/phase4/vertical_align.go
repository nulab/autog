package phase4

import "github.com/nulab/autog/internal/graph"

func execVerticalAlign(g *graph.DGraph, params graph.Params) {
	maxW := 0.0
	for _, layer := range g.Layers {
		layer.H = 0.0
		layer.W = 0.0
		for i, n := range layer.Nodes {
			layer.W += n.W
			// add node spacing except after the last node
			if i < len(layer.Nodes)-1 {
				layer.W += params.NodeSpacing
			}
			layer.H = max(layer.H, n.H)
		}
		// find largest layer
		maxW = max(maxW, layer.W)
	}

	for _, layer := range g.Layers {
		pos := (maxW - layer.W) / 2
		for _, n := range layer.Nodes {
			n.X = pos
			pos += n.W + params.NodeSpacing
		}
	}
}

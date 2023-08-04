package phase4

import "github.com/nulab/autog/graph"

func execVerticalAlign(g *graph.DGraph, params graph.Params) {
	maxW := 0.0
	for _, layer := range g.Layers {
		layer.H = 0.0
		layer.W = 0.0
		var last *graph.Node
		for _, n := range layer.Nodes {
			if last != nil {
				layer.W += params.NodeSpacing
			}
			layer.W += params.NodeMargin*2 + n.W
			layer.H = max(layer.H, n.H)
			last = n
		}
		maxW = max(maxW, layer.W)
	}

	for _, layer := range g.Layers {
		pos := (maxW - layer.W) / 2
		var last *graph.Node
		for _, n := range layer.Nodes {
			if last != nil {
				pos += params.NodeSpacing
			}
			pos += params.NodeMargin
			n.X = pos
			pos += n.W + params.NodeMargin
			last = n
		}
	}
}

package positioning

import "github.com/nulab/autog/graph"

const (
	defaultNodeMargin  = 20.0
	defaultNodeSpacing = 40.0
)

func execVerticalAlign(g *graph.DGraph) {
	maxW := 0.0
	for _, layer := range g.Layers {
		layer.H = 0.0
		layer.W = 0.0
		var last *graph.Node
		for _, n := range layer.Nodes {
			if last != nil {
				layer.W += defaultNodeSpacing
			}
			layer.W += defaultNodeMargin*2 + n.W
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
				pos += defaultNodeSpacing
			}
			pos += defaultNodeMargin
			n.X = pos
			pos += n.W + defaultNodeMargin
			last = n
		}
	}
}

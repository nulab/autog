package phase2

import (
	"github.com/nulab/autog/internal/graph"
)

func execLongestPath(g *graph.DGraph) {
	for _, n := range g.Nodes {
		n.Layer = -1
	}
	maxh := 0
	for _, n := range g.Nodes {
		h := followLongestPath(n)
		maxh = max(maxh, h)
	}
}

func followLongestPath(n *graph.Node) int {
	if n.Layer >= 0 {
		return n.Layer
	}
	maxh := 0
	for _, e := range n.In {
		if e.SelfLoops() {
			continue
		}

		h := followLongestPath(e.ConnectedNode(n))
		maxh = max(maxh, h+1)
	}

	n.Layer = maxh
	return maxh
}

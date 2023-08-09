package phase2

import (
	"github.com/nulab/autog/graph"
)

func execLongestPath(g *graph.DGraph) {
	height := graph.NodeIntMap{}

	for _, n := range g.Nodes {
		height[n] = -1
	}

	nlayers := 0
	for _, n := range g.Nodes {
		followLongestPath(n, height, &nlayers)
	}
}

func followLongestPath(n *graph.Node, height graph.NodeIntMap, nlayers *int) int {
	if height[n] >= 0 {
		return height[n]
	}
	nodeh := 1
	// sinks have no out-edges, so will yield 1
	for _, e := range n.Out {
		if e.SelfLoops() {
			continue
		}

		h := followLongestPath(e.ConnectedNode(n), height, nlayers)
		nodeh = max(nodeh, h+e.Delta)
	}
	*nlayers = max(*nlayers, nodeh)
	n.Layer = *nlayers - nodeh
	height[n] = nodeh
	return nodeh
}

func fillLayers(g *graph.DGraph) {
	highest := 0
	for i := range g.Layers {
		highest = max(highest, i)
	}
	for i := 0; i < highest; i++ {
		_, ok := g.Layers[i]
		if !ok {
			g.Layers[i] = &graph.Layer{Index: i}
		}
	}
}

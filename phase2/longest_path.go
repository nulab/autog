package phase2

import (
	"sort"

	"github.com/nulab/autog/graph"
)

func execLongestPath(g *graph.DGraph) {
	height := graph.NodeIntMap{}

	for _, n := range g.Nodes {
		height[n] = -1
	}
	nodes := make([]*graph.Node, len(g.Nodes))
	copy(nodes, g.Nodes)

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Outdeg() > nodes[j].Outdeg() ||
			(nodes[i].Outdeg() == nodes[j].Outdeg() && nodes[i].Indeg() < nodes[j].Indeg())
	})

	nlayers := 0
	for _, n := range nodes {
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

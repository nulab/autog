package connected

import (
	"maps"

	ig "github.com/nulab/autog/internal/graph"
)

func Components(g *ig.DGraph) []*ig.DGraph {
	visitedN := make(ig.NodeSet)
	visitedE := make(ig.EdgeSet)
	walkDfs(g.Nodes[0], visitedN, visitedE)

	// if all nodes were visited at the first dfs
	// then there is only one connected component and that is G itself
	if len(visitedN) == len(g.Nodes) {
		return []*ig.DGraph{g}
	}

	cnncmp := make([]*ig.DGraph, 0, 2) // this has at least 2 connected components
	cnncmp = append(cnncmp, &ig.DGraph{Nodes: visitedN.Keys(), Edges: visitedE.Keys()})

	for _, n := range g.Nodes {
		if !visitedN[n] {
			ns := make(ig.NodeSet)
			es := make(ig.EdgeSet)
			walkDfs(n, ns, es)
			cnncmp = append(cnncmp, &ig.DGraph{Nodes: ns.Keys(), Edges: es.Keys()})
			maps.Copy(visitedN, ns)
		}
	}
	return cnncmp
}

func walkDfs(n *ig.Node, visitedN ig.NodeSet, visitedE ig.EdgeSet) {
	visitedN[n] = true
	n.VisitEdges(func(e *ig.Edge) {
		if !visitedE[e] {
			visitedE[e] = true
			walkDfs(e.ConnectedNode(n), visitedN, visitedE)
		}
	})
}

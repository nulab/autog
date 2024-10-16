package graph

import (
	"maps"
)

func (g *DGraph) ConnectedComponents() []*DGraph {
	visitedN := make(NodeSet)
	visitedE := make(EdgeSet)
	walkDfs(g.Nodes[0], visitedN, visitedE)

	// if all nodes were visited at the first dfs
	// then there is only one connected component and that is G itself
	if len(visitedN) == len(g.Nodes) {
		return []*DGraph{g}
	}

	cnncmp := make([]*DGraph, 0, 2) // this has at least 2 connected components
	cnncmp = append(cnncmp, &DGraph{Nodes: visitedN.Keys(), Edges: visitedE.Keys()})

	for _, n := range g.Nodes {
		if !visitedN[n] {
			ns := make(NodeSet)
			es := make(EdgeSet)
			walkDfs(n, ns, es)
			cnncmp = append(cnncmp, &DGraph{Nodes: ns.Keys(), Edges: es.Keys()})
			maps.Copy(visitedN, ns)
		}
	}
	return cnncmp
}

func walkDfs(n *Node, visitedN NodeSet, visitedE EdgeSet) {
	visitedN[n] = true
	n.VisitEdges(func(e *Edge) {
		if !visitedE[e] {
			visitedE[e] = true
			walkDfs(e.ConnectedNode(n), visitedN, visitedE)
		}
	})
}

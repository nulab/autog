package phase1

import "github.com/nulab/autog/internal/graph"

func hasCycles(g *graph.DGraph) bool {
	visited := graph.NodeSet{}
	finished := graph.NodeSet{}
	for _, n := range g.Nodes {
		if !visited[n] && !finished[n] {
			if visit(n, visited, finished) {
				return true
			}
		}
	}
	return false
}

func visit(n *graph.Node, visited, finished graph.NodeSet) bool {
	visited[n] = true
	for _, e := range n.Out {
		if e.SelfLoops() {
			continue // ignore self loops
		}
		m := e.ConnectedNode(n)
		if visited[m] {
			return true
		}
		if !finished[m] {
			if visit(m, visited, finished) {
				return true
			}

		}
	}
	visited[n] = false
	finished[n] = true
	return false
}

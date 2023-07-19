package graph

func (g *DGraph) ConnectedComponents() []*DGraph {
	var subgs []*DGraph
	visited := NodeSet{}
	for _, n := range g.Nodes {
		c := connectedSubgraph(n, visited)
		if c != nil {
			subgs = append(subgs, &DGraph{Nodes: c})
		}
	}
	return subgs
}

func connectedSubgraph(n *Node, visited NodeSet) []*Node {
	if visited[n] {
		return nil
	}
	visited[n] = true

	subg := []*Node{n}
	for itr := n.EdgeIter(); itr.HasNext(); {
		e := itr.Next()
		if e.ConnectedNode(n) == n {
			continue // self-loop
		}
		ns := connectedSubgraph(e.ConnectedNode(n), visited)
		subg = append(subg, ns...)
	}
	return subg
}

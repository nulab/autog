package graph

func (g *DGraph) ConnectedComponents() []*DGraph {
	var subgs []*DGraph
	visited := NodeSet{}
	for _, n := range g.Nodes {
		c := connectedSubgraph(n, visited)
		if c != nil {

			subgs = append(subgs, &DGraph{Nodes: c, Edges: edgesOf(c)})
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

func edgesOf(nodes []*Node) []*Edge {
	res := []*Edge{}
	visited := EdgeSet{}
	for _, n := range nodes {
		res = append(res, collectEdges(n, visited)...)
	}
	return res
}

func collectEdges(n *Node, visited EdgeSet) []*Edge {
	var res []*Edge
	for itr := n.EdgeIter(); itr.HasNext(); {
		e := itr.Next()
		if visited[e] {
			continue
		}
		visited[e] = true
		res = append(res, e)
	}
	return res
}

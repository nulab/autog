package graph

import (
	"strings"
)

// todo maybe not needed
type EdgeList []*Edge

func (s EdgeList) String() string {
	bld := strings.Builder{}
	for _, e := range s {
		bld.WriteString(e.String())
		bld.WriteRune('\n')
	}
	return bld.String()
}

// SpanningTree computes a spanning tree t of g, marks edges that belong to t as such.
// It returns a list of all edges.
func (g *DGraph) SpanningTree() EdgeList {
	root := g.Nodes[0]

	visitedNodes := NodeSet{}
	visitedNodes[root] = true

	visitedEdges := EdgeSet{}

	queue := []*Node{}
	queue = append(queue, root)

	var sEdges []*Edge
	var n *Node

	for len(queue) > 0 {
		n, queue = queue[0], queue[1:]

		nEdges := n.Edges()
		for _, e := range nEdges {
			w := e.ConnectedNode(n)
			if !visitedNodes[w] {
				visitedNodes[w] = true
				visitedEdges[e] = true

				e.IsInSpanningTree = true
				sEdges = append(sEdges, e)
				queue = append(queue, w)
			} else {
				if !visitedEdges[e] {
					visitedEdges[e] = true
					sEdges = append(sEdges, e)
				}
			}
		}
	}
	return sEdges
}

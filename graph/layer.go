package graph

type Layer struct {
	Nodes []*Node
	Index int
	Size
}

// CountCrossings returns the number of edge crossings between this layer and the layer above.
func (layer *Layer) CountCrossings() int {
	if layer.Index == 0 {
		return 0
	}
	edges := []*Edge{}
	for _, n := range layer.Nodes {
		visit, next := n.AllEdges()
		for next {
			next = visit(func(e *Edge) {
				if e.ConnectedNode(n).Layer < n.Layer {
					edges = append(edges, e)
				}
			})
		}
	}
	return countCrossings(edges)
}

// comparable 2-tuple of edge pointers
type pair [2]*Edge

// naive O(n^2) count of crossings between edges in a set
func countCrossings(edges []*Edge) int {
	crossings := 0
	visited := map[pair]bool{}
	for _, e := range edges {
		for _, f := range edges {
			if f == e || visited[pair{e, f}] || visited[pair{f, e}] {
				continue
			}
			visited[pair{e, f}] = true
			if e.Crosses(f) {
				crossings++
			}
		}
	}
	return crossings
}

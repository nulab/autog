package phase3

import (
	"github.com/nulab/autog/graph"
)

// Barth-Mutzel crossing counting algorithm that runs in O(|E|*log|V|) described in
// "Simple and Efficient Bilayer Cross Counting", Journal of Graph Algorithms and Applications, 26 August 2002
// https://pdfs.semanticscholar.org/272d/73edce86bcfac3c82945042cf6733ad281a0.pdf
func countCrossings(l1, l2 *graph.Layer) int {
	upper, lower := orderedLayers(l1, l2)
	edges := inLayerEdges(upper, lower)

	// obtain a radix-sorted slice of target nodes
	nodes := radixsort(upper, lower, edges)

	// initialize the tree size as a power of 2 that is just large enough to hold all vertices of the smaller layer
	// in its leafs; the paper uses the inequality 2^(c-1) < q <= 2^c where q is the size of the smaller layer
	q := min(l1.Len(), l2.Len())
	k := 1
	for k < q {
		k *= 2
	}
	size := 2*k - 1
	tree := make([]int, size)

	k -= 1 // now k is the index at which tree leafs start

	crosscount := 0
	for _, n := range nodes {
		i := n.LayerPos + k
		tree[i]++
		for i > 0 {
			// sum value of the right tree node when encountering a left tree node
			if i%2 != 0 {
				crosscount += tree[i+1]
			}
			// move up to the parent
			i = (i - 1) / 2
			tree[i]++
		}
	}
	return crosscount
}

// sorts the edges lexicographically based on the position of their upper and lower vertices,
// basically e(u,v) < f(w,x) iff (pos(u) < pos(w) || (pos(u) == pos(w) && pos{v) < pos(x))),
// then collects target nodes in that order; runs in O(|E|+|V1||V2|)
func radixsort(upper, lower *graph.Layer, es []*graph.Edge) []*graph.Node {
	m := upper.Len()
	n := lower.Len()
	// using nodes because it gives a meaningful zero value
	mat := make([][]*graph.Node, m)
	for i := range mat {
		mat[i] = make([]*graph.Node, n)
	}
	for _, e := range es {
		src, tgt := orderedEdgeNodes(upper.Index, e)
		mat[src.LayerPos][tgt.LayerPos] = tgt
	}
	nodes := make([]*graph.Node, len(es))
	k := 0
	for i := range mat {
		for j := range mat[i] {
			if mat[i][j] != nil {
				nodes[k] = mat[i][j]
				k++
			}
		}
	}
	return nodes
}

// collect edges between layers l1 and l2; runs in O(|E|)
func inLayerEdges(upper, lower *graph.Layer) []*graph.Edge {
	es := []*graph.Edge{}
	for _, n := range upper.Nodes {
		n.VisitEdges(func(e *graph.Edge) {
			// compare regardless of order
			if bit(e.From.Layer, e.To.Layer) == bit(upper.Index, lower.Index) {
				es = append(es, e)
			}
		})

	}
	return es
}

func bit(a, b int) uint64 {
	return (1 << a) | (1 << b)
}

// returns the layers as a tuple ordered by number of nodes
func orderedLayers(l1, l2 *graph.Layer) (upper, lower *graph.Layer) {
	if l1.Len() > l2.Len() {
		return l1, l2
	}
	return l2, l1
}

// returns the nodes of e as if e were directed from the biggest to the smallest layer
func orderedEdgeNodes(upperl int, e *graph.Edge) (upper, lower *graph.Node) {
	if e.From.Layer == upperl {
		return e.From, e.To
	}
	return e.To, e.From
}

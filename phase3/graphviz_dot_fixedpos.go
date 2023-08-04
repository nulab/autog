package phase3

import "github.com/nulab/autog/graph"

type fixedPositions struct {
	mustAfter  graph.NodeMap
	mustBefore graph.NodeMap
}

func initFixedPositions(edges []*graph.Edge) fixedPositions {
	mustAfter := graph.NodeMap{}
	mustBefore := graph.NodeMap{}
	for _, e := range edges {
		if e.From.Layer == e.To.Layer {
			mustAfter[e.To] = e.From
			mustBefore[e.From] = e.To
		}
	}
	return fixedPositions{mustAfter, mustBefore}
}

// head returns the first element in a same-layer transitive closure to which k belongs, and the number of edges
// that separate k and the head;
// or returns k itself and 0 if k doesn't belong to any such closure
func (fp *fixedPositions) head(k *graph.Node) (*graph.Node, int) {
	v := k
	i := 0
	for n, ok := fp.mustAfter[k]; ok; n, ok = fp.mustAfter[n] {
		v = n
		i++
	}
	return v, i
}

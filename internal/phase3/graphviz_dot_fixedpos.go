package phase3

import (
	"github.com/nulab/autog/internal/graph"
)

type fixedPositions struct {
	mustAfter  graph.NodeMap
	mustBefore graph.NodeMap
}

func initFixedPositions(edges []*graph.Edge) fixedPositions {
	type lln struct {
		n          *graph.Node
		prev, next *lln
	}
	known := map[*graph.Node]*lln{}
	chains := []*lln{}

	for _, e := range edges {
		if e.IsFlat() {
			switch {
			case known[e.From] == nil && known[e.To] == nil:
				_1 := &lln{n: e.From}
				_2 := &lln{n: e.To}
				_1.next = _2
				_2.prev = _1
				known[e.From] = _1
				known[e.To] = _2
				chains = append(chains, _1)

			case known[e.From] != nil && known[e.To] == nil:
				_n := known[e.From]
				for _n.next != nil {
					_n = _n.next
				}
				_2 := &lln{n: e.To}
				_n.next = _2
				_2.prev = _n
				known[e.To] = _2

			case known[e.From] == nil && known[e.To] != nil:
				_n := known[e.To]
				_1 := &lln{n: e.From}
				if _n.prev != nil {
					_1.prev = _n.prev
					_n.prev.next = _1
					_n.prev = _1
					_1.next = _n
				} else {
					for i := range chains {
						if chains[i] == _n {
							chains[i] = _1
						}
					}
					_n.prev = _1
					_1.next = _n

				}
				known[e.From] = _1
			}
		}
	}

	mustAfter := graph.NodeMap{}
	mustBefore := graph.NodeMap{}
	for _, a := range chains {
		for a.next != nil {
			mustBefore[a.n] = a.next.n
			a = a.next
			mustAfter[a.n] = a.prev.n
		}
	}
	return fixedPositions{mustAfter, mustBefore}
}

func walkFlat(n *graph.Node, visited graph.EdgeSet) {
	for _, f := range n.Out {
		if visited[f] {
			continue
		}
	}
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

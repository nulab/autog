package phase5

import (
	"github.com/nulab/autog/graph"
)

type Alg uint8

const (
	NoRouting Alg = iota
	_endAlg
)

func (alg Alg) IsValid() bool {
	return alg < _endAlg
}

func (alg Alg) Process(g *graph.DGraph, params graph.Params) {
	// remove virtual nodes
	i := 0
loop:
	for i < len(g.Nodes) {
		n := g.Nodes[i]
		i++
		if n.IsVirtual {
			if n.Indeg() != 1 || n.Outdeg() != 1 {
				panic("virtual node has not exactly 1 in-edge and 1 out-edge")
			}
			in, out := n.In[0], n.Out[0]
			// 	 from-in-to  from-out-to
			// n1 --------> n --------> n2
			in.To = out.To
			// then replace 'out' with 'in' in n2's incoming edge list
			n2 := out.To
			g.Nodes = append(g.Nodes[:i-1], g.Nodes[i:]...)
			for j, f := range n2.In {
				if f == out {
					n2.In[j] = in
					for k, r := range g.Edges {
						if r == out {
							g.Edges = append(g.Edges[:k], g.Edges[k+1:]...)
						}
					}
					goto loop
				}
			}
		}
	}

	switch alg {
	case NoRouting:
		return
	default:
		panic("routing: unknown alg value")
	}
}

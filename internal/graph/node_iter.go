package graph

import "iter"

func (n *Node) VisitEdges(visit func(*Edge)) {
	for e := range n.allEdges() {
		visit(e)
	}
}

func (n *Node) allEdges() iter.Seq[*Edge] {
	return func(yield func(*Edge) bool) {
		for i, b := 0, true; i < n.Deg(); i++ {
			if i < n.Indeg() {
				b = yield(n.In[i])
			} else {
				b = yield(n.Out[i-n.Indeg()])
			}
			if !b {
				return
			}
		}
	}
}

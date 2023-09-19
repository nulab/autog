package graph

import (
	"sync"
)

type Edge struct {
	_ [0]sync.Mutex
	edge
}

// todo: collapse flags into one bitmask field
type edge struct {
	From, To         *Node
	Delta            int          // edge's minimum length, used in network simplex layerer
	Weight           int          // edge's weight, used in network simplex layerer
	IsInSpanningTree bool         // whether the edge is in the graph spanning tree, used in network simplex layerer
	IsReversed       bool         // whether the edge has been reversed to break cycles, set in phase 1
	CutValue         int          // edge's cut value, used in network simplex layerer
	Points           [][2]float64 // edge's control points, used in phase 5
	ArrowHeadStart   bool         // arrowhead position, 0: start, 1: end
}

func NewEdge(from, to *Node, weight int) *Edge {
	return &Edge{
		edge: edge{
			From:   from,
			To:     to,
			Delta:  1,
			Weight: weight,
		},
	}
}

func (e *Edge) SelfLoops() bool {
	return e.From == e.To
}

func (e *Edge) IsFlat() bool {
	return e.From.Layer == e.To.Layer
}

func (e *Edge) Reverse() {
	from, to := e.From, e.To

	from.Out.Remove(e)
	to.In.Remove(e)

	from.In.Add(e)
	to.Out.Add(e)

	e.From = to
	e.To = from
	e.IsReversed = !e.IsReversed
}

func (e *Edge) ConnectedNode(n *Node) *Node {
	if e.To != n {
		return e.To
	}
	return e.From
}

func (e *Edge) Crosses(f *Edge) bool {
	etop, ebtm := e.From, e.To
	if e.To.Layer > e.From.Layer {
		etop, ebtm = e.To, e.From
	}
	ftop, fbtm := f.From, f.To
	if f.To.Layer > f.From.Layer {
		ftop, fbtm = f.To, f.From
	}
	return (etop.LayerPos < ftop.LayerPos && ebtm.LayerPos > fbtm.LayerPos) ||
		(etop.LayerPos > ftop.LayerPos && ebtm.LayerPos < fbtm.LayerPos)
}

// Type returns the edge type as follows:
//   - 0: both of e's adjacent nodes are concrete nodes
//   - 1: exactly one of e's adjacent nodes is virtual
//   - 2: both of e's adjacent nodes are virtual
func (e *Edge) Type() int {
	// todo: might return an enum instead
	if !e.From.IsVirtual && !e.To.IsVirtual {
		return 0
	}
	if e.From.IsVirtual != e.To.IsVirtual {
		return 1
	}
	if e.From.IsVirtual && e.To.IsVirtual {
		return 2
	}
	panic("edge type cases aren't exhaustive")
}

func (e *Edge) String() string {
	s := e.From.ID + " -> " + e.To.ID
	if e.IsReversed {
		s += " (rev)"
	}
	if !e.IsInSpanningTree {
		s += " (non-stree)"
	}
	return s
}

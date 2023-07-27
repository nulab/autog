package graph

import "github.com/nulab/autog/internal/pragma"

type Edge struct {
	pragma.NotCopiable

	ID       string
	From, To *Node

	// variables not relevant to the final output
	IsInSpanningTree bool
	IsReversed       bool
	CutValue         int
}

func (e *Edge) SelfLoops() bool {
	return e.From == e.To
}

func (e *Edge) Reverse() {
	n := e.From
	e.From = e.To
	e.To = n
	e.IsReversed = !e.IsReversed
}

func (e *Edge) ConnectedNode(n *Node) *Node {
	if e.To != n {
		return e.To
	}
	return e.From
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

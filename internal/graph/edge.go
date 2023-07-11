package graph

type Edge struct {
	// ID       int
	From, To *Node
}

func (e *Edge) Reverse() {
	n := e.From
	e.From = e.To
	e.To = n
}

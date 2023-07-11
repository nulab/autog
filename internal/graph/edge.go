package graph

type Edge struct {
	ID       string
	From, To *Node
	reversed bool
}

func (e *Edge) Reverse() {
	n := e.From
	e.From = e.To
	e.To = n
	e.reversed = true
}

func (e *Edge) IsReversed() bool {
	return e.reversed
}

func (e Edge) String() string {
	return e.From.ID + " -> " + e.To.ID
}

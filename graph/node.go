package graph

import (
	"sync"
)

type Node struct {
	_ [0]sync.Mutex

	ID        string
	In, Out   EdgeList
	Layer     int
	LayerPos  int
	IsVirtual bool
	Size
}

func (n *Node) String() string {
	return n.ID
}

// Indeg returns the number of incoming edges
func (n *Node) Indeg() int {
	return len(n.In)
}

// Outdeg returns the number of outgoing edges
func (n *Node) Outdeg() int {
	return len(n.Out)
}

// Deg returns the total number of incoming and outgoing edges
func (n *Node) Deg() int {
	return n.Indeg() + n.Outdeg()
}

package graph

import (
	"fmt"
	"sync"
)

type Node struct {
	_ [0]sync.Mutex

	ID        string
	In, Out   EdgeList
	Layer     int
	LayerPos  int
	IsVirtual bool

	// Size stores the dimensions of the node and its computed (x,y) coordinates.
	// The point (x,y) is the top-left corner of the node, and (0,0) is the top-left corner of the drawing plane.
	Size
}

func (n *Node) String() string {
	return n.ID
}

func (n *Node) SVG() string {
	return fmt.Sprintf(
		`<rect id="%s" x="%f" y="%f" width="%f" height="%f" style="fill: none; stroke: blue;" />`,
		"autog-node-"+n.ID,
		n.X, n.Y, n.W, n.H,
	)
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

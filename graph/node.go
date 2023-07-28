package graph

import "github.com/nulab/autog/internal/pragma"

type Node struct {
	pragma.NotCopiable

	ID        string
	In, Out   []*Edge // this could also be list of nodes
	Layer     int     // todo: this probably shouldn't be visible outside
	LayerPos  int
	IsVirtual bool
	Size
}

func (n *Node) String() string {
	return n.ID
}

// Number of incoming edges
func (n *Node) Indeg() int {
	return len(n.In)
}

// Number of outgoing edges
func (n *Node) Outdeg() int {
	return len(n.Out)
}

// Total number of incoming and outgoing edges
func (n *Node) Deg() int {
	return n.Indeg() + n.Outdeg()
}

func (n *Node) Edges() (edges []*Edge) {
	edges = make([]*Edge, 0, n.Deg())
	edges = append(edges, n.In...)
	edges = append(edges, n.Out...)
	return
}

func (n *Node) AdjacentNodes() (nodes []*Node) {
	// collect target nodes of outgoing edges (the source node is always 'n' itself)
	for _, outE := range n.Out {
		nodes = append(nodes, outE.To)
	}
	return
}

package graph

import "github.com/nulab/autog/internal/graph"

type Layout struct {
	Nodes []Node
	Edges []Edge
}

type Node struct {
	ID string
	graph.Size
}

type Edge struct {
	ID string
	// todo: points to draw the edge
}

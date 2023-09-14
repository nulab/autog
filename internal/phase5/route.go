package phase5

import "github.com/nulab/autog/internal/graph"

// holds information about the route an edge has to go through
type route struct {
	// nodes that the edges passes through
	// the start node is guaranteed to be the node with lesser layer, or lesser layer order
	// the end node is guaranteed to be the node with greater layer, or greater layer order
	ns []*graph.Node
}

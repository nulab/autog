package autog

import (
	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/phase1"
)

// todo: to decrease coupling between client code and the graph types used here, the layout could take as params
// a simple adjacency list and a struct or map with the node properties (width, height, etc.)
// then the graph package could become internal
func Layout(graph *graph.DGraph, opts ...option) *graph.DGraph {
	layoutOpts := defaultOptions
	for _, opt := range opts {
		opt(&layoutOpts)
	}
	defer layoutOpts.params.Monitor.Close()

	pipeline := [...]processor{
		layoutOpts.p1, // cycle breaking
		layoutOpts.p2, // layering
		layoutOpts.p3, // ordering
		layoutOpts.p4, // positioning
		layoutOpts.p5, // edge routing
		phase1.RestoreEdges,
	}

	for _, phase := range pipeline {
		phase.Process(graph, layoutOpts.params)
	}

	return graph
}

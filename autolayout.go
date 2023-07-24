package autog

import (
	"github.com/vibridi/autog/graph"
)

func Layout(graph *graph.DGraph, opts ...option) *graph.DGraph {
	layoutOpts := defaultOptions
	for _, opt := range opts {
		opt(&layoutOpts)
	}

	pipeline := [...]processor{
		layoutOpts.p1, // cycle breaking
		layoutOpts.p2, // layering
		// node ordering (input: layered graph, output: layered graph with node order)
		// node positioning (input: ???, output: layered graph with node coordinates)
		// edge routing
		// todo: restore reverted edges if necessary
	}

	for _, phase := range pipeline {
		phase.Process(graph)
	}

	return graph
}

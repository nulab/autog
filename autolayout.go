package autog

import (
	"github.com/vibridi/autog/graph"
)

type nodeCoords int8 // placeholder type

// todo: if DGraph is internal, nobody can actually construct instances
func Layout(graph *graph.DGraph, opts ...option) nodeCoords {
	layoutOpts := defaultOptions
	for _, opt := range opts {
		opt(&layoutOpts)
	}

	pipeline := [...]processor{
		layoutOpts.p1, // cycle break (input: directed graph, output: directed acyclic graph)
		layoutOpts.p2,
		// todo: restore reverted edges if necessary
	}

	for _, phase := range pipeline {
		phase.Process(graph)
	}

	// layering (input: DAG, output: layered graph)

	// node ordering (input: layered graph, output: layered graph with node order)

	// node positioning (input: ???, output: layered graph with node coordinates)

	// edge routing

	// layered graph
	return 0
}

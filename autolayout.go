package autog

import (
	"github.com/vibridi/autog/internal/cyclebreaking"
	"github.com/vibridi/autog/internal/graph"
)

type nodeCoords int8 // placeholder type

// todo: if DGraph is internal, nobody can actually construct instances
func Layout(graph *graph.DGraph) nodeCoords {

	pipeline := [5]phase{
		cyclebreaking.Greedy,
		// todo: restore reverted edges if necessary
	}

	for _, p := range pipeline {
		p.Process(graph)
		p.Cleanup()
	}

	// cycle break (input: directed graph, output: acyclic)

	// layering (input: DAG, output: layered graph)

	// node ordering (input: layered graph, output: layered graph with node order)

	// node positioning (input: ???, output: layered graph with node coordinates)

	// edge routing

	// layered graph
	return 0
}

package autog

import (
	"github.com/nulab/autog/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
	"github.com/nulab/autog/internal/processor"
	"github.com/nulab/autog/phase1"
)

// todo: to decrease coupling between client code and the graph types used here, the layout could take as params
// a simple adjacency list and a struct or map with the node properties (width, height, etc.)
// then the graph package could become internal

// todo: add interactive layout

func Layout(graph *graph.DGraph, opts ...Option) *graph.DGraph {
	layoutOpts := defaultOptions
	for _, opt := range opts {
		opt(&layoutOpts)
	}

	imonitor.Set(layoutOpts.monitor)
	defer imonitor.Reset()

	pipeline := [...]processor.P{
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

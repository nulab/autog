package autog

import (
	"github.com/nulab/autog/graph"
	igraph "github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
	"github.com/nulab/autog/internal/processor"
)

// todo: add interactive layout

func Layout(source graph.Source, opts ...Option) *igraph.DGraph {
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
	}

	// obtain the graph struct
	g := source.Generate()

	// run it through the pipeline
	for _, phase := range pipeline {
		phase.Process(g, layoutOpts.params)
	}

	return g
}

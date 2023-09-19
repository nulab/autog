package autog

import (
	"github.com/nulab/autog/graph"
	ig "github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
	"github.com/nulab/autog/internal/processor"
)

// todo: add interactive layout

func Layout(source graph.Source, opts ...Option) graph.Layout {
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

	// populate the graph struct from the graph source
	g := &ig.DGraph{}
	source.Populate(g)

	if layoutOpts.params.NodeFixedSizeFunc != nil {
		for _, n := range g.Nodes {
			layoutOpts.params.NodeFixedSizeFunc(n)
		}
	}

	// run it through the pipeline
	for _, phase := range pipeline {
		phase.Process(g, layoutOpts.params)
	}

	// return only relevant data to the caller
	out := graph.Layout{
		Nodes: make([]graph.Node, len(g.Nodes)),
		Edges: make([]graph.Edge, len(g.Edges)),
	}
	for i, n := range g.Nodes {
		if n.IsVirtual {
			continue
		}
		out.Nodes[i] = graph.Node{ID: n.ID, Size: n.Size}
	}
	for i, e := range g.Edges {
		out.Edges[i] = graph.Edge{
			Points:         e.Points,
			ArrowHeadStart: e.ArrowHeadStart,
		}
	}
	return out
}

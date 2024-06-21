package autog

import (
	"slices"

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
	g := from(source)

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
		Nodes: make([]graph.Node, 0, len(g.Nodes)),
		Edges: make([]graph.Edge, 0, len(g.Edges)),
	}
	for _, n := range g.Nodes {
		if n.IsVirtual && !layoutOpts.output.keepVirtualNodes {
			continue
		}
		out.Nodes = append(out.Nodes, graph.Node{
			ID:   n.ID,
			Size: n.Size,
		})
		// todo: clients can't reliably tell virtual nodes from concrete nodes
	}
	if !layoutOpts.output.keepVirtualNodes {
		out.Nodes = slices.Clip(out.Nodes)
	}

	for _, e := range g.Edges {
		out.Edges = append(out.Edges, graph.Edge{
			FromID:         e.From.ID,
			ToID:           e.To.ID,
			Points:         e.Points,
			ArrowHeadStart: e.ArrowHeadStart,
		})
	}
	return out
}

func from(source graph.Source) *ig.DGraph {
	switch t := source.(type) {
	case *ig.DGraph:
		// special case for when the graph source is already a DGraph
		// this happens only during unit testing
		return t
	default:
		g := &ig.DGraph{}
		source.Populate(g)
		return g
	}
}

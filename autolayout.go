package autog

import (
	"slices"

	"github.com/nulab/autog/graph"
	ig "github.com/nulab/autog/internal/graph"
	"github.com/nulab/autog/internal/graph/connected"
	imonitor "github.com/nulab/autog/internal/monitor"
	"github.com/nulab/autog/internal/processor"
	"github.com/nulab/autog/internal/processor/preprocessor"
)

// Layout executes the layout algorithm on the graph G obtained from source. It panics if G contains no nodes.
func Layout(source graph.Source, opts ...Option) graph.Layout {
	layoutOpts := defaultOptions
	for _, opt := range opts {
		opt(&layoutOpts)
	}

	imonitor.Set(layoutOpts.monitor)
	defer imonitor.Reset()

	pipeline := []processor.P{
		layoutOpts.p1, // cycle breaking
		layoutOpts.p2, // layering
		layoutOpts.p3, // ordering
		layoutOpts.p4, // positioning
		layoutOpts.p5, // edge routing
	}

	// populate the graph struct from the graph source
	G := from(source)

	if len(G.Nodes) == 0 {
		panic("autog: node set is empty")
	}

	if layoutOpts.params.NodeFixedSizeFunc != nil {
		for _, n := range G.Nodes {
			layoutOpts.params.NodeFixedSizeFunc(n)
		}
	}

	// return only relevant data to the caller
	out := graph.Layout{}

	// shift disconnected sub-graphs to the right
	shift := 0.0

	// process each connected components and collect results into the same layout output
	for _, g := range connected.Components(G) {
		if len(g.Nodes) == 0 {
			panic("autog: connected sub-graph node set is empty: this might be a bug")
		}

		out.Nodes = slices.Grow(out.Nodes, len(g.Nodes))
		out.Edges = slices.Grow(out.Edges, len(g.Edges))

		// pre-processing
		restoreSelfLoops := preprocessor.IgnoreSelfLoops(g)

		// run subgraph through the pipeline
		for _, phase := range pipeline {
			phase.Process(g, layoutOpts.params)
		}

		// post-processing
		restoreSelfLoops(g)

		// collect nodes
		for _, n := range g.Nodes {
			if n.IsVirtual && !layoutOpts.output.keepVirtualNodes {
				continue
			}

			m := graph.Node{
				ID:   n.ID,
				Size: n.Size,
			}
			// apply subgraph's left shift
			m.X += shift

			out.Nodes = append(out.Nodes, m)
			// todo: clients can't reliably tell virtual nodes from concrete nodes
		}

		// collect edges
		for _, e := range g.Edges {
			f := graph.Edge{
				FromID:         e.From.ID,
				ToID:           e.To.ID,
				Points:         slices.Clone(e.Points),
				ArrowHeadStart: e.ArrowHeadStart,
			}
			// apply subgraph's left shift
			for i := range f.Points {
				f.Points[i][0] += shift
			}

			out.Edges = append(out.Edges, f)
		}

		// compute shift for subsequent subgraphs
		rightmostX := 0.0
		for _, l := range g.Layers {
			if len(l.Nodes) == 0 {
				// empty layers don't affect the shift
				continue
			}
			n := l.Nodes[len(l.Nodes)-1]
			rightmostX = max(rightmostX, n.X+n.W)
		}
		shift += rightmostX + layoutOpts.params.NodeSpacing
	}

	if !layoutOpts.output.keepVirtualNodes {
		out.Nodes = slices.Clip(out.Nodes)
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

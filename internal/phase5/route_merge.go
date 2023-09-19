package phase5

import (
	"slices"

	"github.com/nulab/autog/internal/graph"
)

const (
	edgeTypeNoneVirtual = iota
	edgeTypeOneVirtual
	edgeTypeBothVirtual
)

// merge long edges and at the same time collect route information;
// merged edges are removed from the graph edge list
func mergeLongEdges(g *graph.DGraph) []routableEdge {
	routes := make([]routableEdge, 0, len(g.Edges))
	for _, e := range g.Edges {
		switch e.Type() {
		case edgeTypeNoneVirtual:
			u, v := orderedNodes(e)
			e.ArrowHeadStart = e.IsReversed
			routes = append(routes, routableEdge{e, route{[]*graph.Node{u, v}}})

		case edgeTypeOneVirtual:
			// process each chain of virtual nodes only in the direction of the edge
			if e.From.IsVirtual {
				continue
			}
			routes = append(routes, routableEdge{e, route{reduceForward(g, e)}})

		case edgeTypeBothVirtual:
			// skip, eventually it will be processed when encountering a type 1 edge
		}
	}
	return routes
}

// returns the nodes adjacent to this edge in lexicographical order (lesser coordinates first)
func orderedNodes(e *graph.Edge) (u, v *graph.Node) {
	switch {
	case e.From.Layer < e.To.Layer:
		return e.From, e.To

	case e.From.Layer > e.To.Layer:
		return e.To, e.From

	default:
		if e.From.LayerPos < e.To.LayerPos {
			return e.From, e.To
		}
		return e.To, e.From
	}
}

// merges a chain of edges connecting virtual nodes into the starting edge e and returns a route
// containing the ordered list of nodes that this edge passes through
func reduceForward(g *graph.DGraph, e *graph.Edge) (ns []*graph.Node) {
	ns = append(ns, e.From)
	for e.To.IsVirtual {
		if len(e.To.Out) != 1 {
			panic("edge routing: virtual node doesn't have exactly one exit edge")
		}
		ns = append(ns, e.To)

		f := e.To.Out[0]
		v := f.To
		v.In.Remove(f)
		v.In.Add(e)
		e.To = v
		g.Edges.Remove(f)
	}
	// append also the non-virtual end node
	ns = append(ns, e.To)
	u, v := orderedNodes(e)
	e.ArrowHeadStart = e.IsReversed
	if ns[0] == v && ns[len(ns)-1] == u {
		slices.Reverse(ns)
	}
	return ns
}

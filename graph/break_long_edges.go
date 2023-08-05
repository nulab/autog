package graph

import "strconv"

func (g *DGraph) BreakLongEdges() {
	for _, e := range g.Edges {
		if e.From.Layer > e.To.Layer {
			panic("edge breaker: upward edge: " + e.From.ID + " -> " + e.To.ID)
		}
	}

	v := 1
	i := 0
loop:
	for i < len(g.Edges) {
		e := g.Edges[i]
		i++
		// target node is below
		if e.To.Layer-e.From.Layer > 1 {
			g.breakLongEdge(e, v)
			v++
			// restart loop
			goto loop
		}
	}
}

func (g *DGraph) breakLongEdge(e *Edge, v int) {
	from, to := e.From, e.To
	// create virtual node
	virtualNode := &Node{
		ID:        "V" + strconv.Itoa(v),
		Layer:     from.Layer + 1,
		IsVirtual: true,
		Size:      Size{H: 100.0, W: 100.0}, // todo: eventually this doesn't belong here
	}
	// set e's target to the virtual node
	e.To = virtualNode
	// add e to virtual node incoming edges
	virtualNode.In = []*Edge{e}
	// create new edge from virtual to e's former target
	f := NewEdge(virtualNode, to, 1)
	f.IsReversed = e.IsReversed
	// add f to virtual node outgoing edges
	virtualNode.Out = []*Edge{f}
	// replace e with f in e's former target incoming edges
	for _, in := range to.In {
		if in == e {
			e = f
		}
	}
	// update the graph's node and edge lists
	g.Edges = append(g.Edges, f)
	g.Nodes = append(g.Nodes, virtualNode)
	g.Layers[virtualNode.Layer].Nodes = append(g.Layers[virtualNode.Layer].Nodes, virtualNode)
}

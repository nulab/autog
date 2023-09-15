package graph

import (
	"strconv"
)

func (g *DGraph) BreakLongEdges() {
	v := 1
	for i := 0; i < len(g.Edges); i++ {
		e := g.Edges[i]
		if e.To.Layer-e.From.Layer > 1 {
			// long edge pointing downward, break regularly
			g.breakLongEdge(e, v)
			v++
		} else if e.From.Layer-e.To.Layer > 1 {
			// long edge pointing upward, temporarily reverse
			e.Reverse()
			e, f := g.breakLongEdge(e, v)
			v++
			// restore direction to both parts
			e.Reverse()
			f.Reverse()
		}
	}
}

func (g *DGraph) breakLongEdge(e *Edge, v int) (*Edge, *Edge) {
	from, to := e.From, e.To
	// create virtual node
	// note that the size of the virtual node may affect positioning algorithms
	// with no size, the node effectively becomes the bend point of the long edge
	virtualNode := &Node{
		ID:        "V" + strconv.Itoa(v),
		Layer:     from.Layer + 1,
		IsVirtual: true,
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
	for i, in := range to.In {
		if in == e {
			to.In[i] = f
			break
		}
	}
	// update the graph's node and edge lists
	g.Edges = append(g.Edges, f)
	g.Nodes = append(g.Nodes, virtualNode)
	g.Layers[virtualNode.Layer].Nodes = append(g.Layers[virtualNode.Layer].Nodes, virtualNode)
	return e, f
}

package phase3

import (
	"strconv"

	"github.com/nulab/autog/internal/graph"
)

func breakLongEdges(g *graph.DGraph) {
	v := 1
	for i := 0; i < len(g.Edges); i++ {
		e := g.Edges[i]
		if e.To.Layer-e.From.Layer > 1 {
			// long edge pointing downward, break regularly
			breakEdge(g, e, v)
			v++
		} else if e.From.Layer-e.To.Layer > 1 {
			// long edge pointing upward, temporarily reverse
			e.Reverse()
			e, f := breakEdge(g, e, v)
			v++
			// restore direction to both parts
			e.Reverse()
			f.Reverse()
		}
	}
}

func breakEdge(g *graph.DGraph, e *graph.Edge, v int) (*graph.Edge, *graph.Edge) {
	from, to := e.From, e.To
	// create virtual node
	// note that the size of the virtual node may affect positioning algorithms
	// with no size, the node effectively becomes the bend point of the long edge
	virtualNode := &graph.Node{
		ID:        "V" + strconv.Itoa(v),
		Layer:     from.Layer + 1,
		IsVirtual: true,
		Size:      graph.Size{W: 20, H: 20},
	}
	// set e's target to the virtual node
	e.To = virtualNode
	// add e to virtual node incoming edges
	virtualNode.In = []*graph.Edge{e}
	// create new edge from virtual to e's former target
	f := graph.NewEdge(virtualNode, to, 1)
	f.IsReversed = e.IsReversed
	// add f to virtual node outgoing edges
	virtualNode.Out = []*graph.Edge{f}
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

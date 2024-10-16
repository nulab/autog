package phase5

import (
	"sort"
	"testing"

	"github.com/nulab/autog/internal/graph"
	"github.com/stretchr/testify/assert"
)

func TestMergeLongEdges(t *testing.T) {
	G := fromEdgeSlice([][]string{
		{"A", "B"},
		{"B", "C"},
		{"B", "D"},
		{"B", "E"},
		{"D", "F"},
		{"A", "V1"}, {"V1", "V2"}, {"V2", "F"},
		{"A", "V3"}, {"V3", "E"},
		{"C", "G"},
		{"B", "V4"}, {"V4", "G"},
		{"G", "A"},
	})
	for _, n := range G.Nodes {
		switch n.ID {
		case "V1", "V2", "V3", "V4":
			n.IsVirtual = true
		}
	}
	for _, e := range G.Edges {
		if isEdge(e, "C", "G") || isEdge(e, "B", "C") || isEdge(e, "B", "V4") || isEdge(e, "V4", "G") {
			e.Reverse()
		}
	}
	G.Layers = make([]*graph.Layer, 5)
	for _, n := range G.Nodes {
		setLayer(G, n)
	}

	assert.Len(t, G.Edges, 14)

	vn := 0
	for _, n := range G.Nodes {
		if n.IsVirtual {
			vn++
		}
	}
	assert.Equal(t, 4, vn)

	routes := mergeLongEdges(G)
	assert.Len(t, G.Edges, 10)

	sort.Slice(routes, func(i, j int) bool {
		if routes[i].ns[0].ID == routes[j].ns[0].ID {
			return routes[i].ns[1].ID < routes[j].ns[1].ID
		}
		return routes[i].ns[0].ID < routes[j].ns[0].ID
	})

	wantRoutesIds := [][]string{
		{"A", "B"},
		{"A", "V1", "V2", "F"},
		{"A", "V3", "E"},
		{"B", "D"},
		{"B", "E"},
		{"C", "B"},
		{"D", "F"},
		{"G", "A"},
		{"G", "A"},
		{"G", "A"},
		{"G", "A"},
		{"G", "C"},
		{"G", "V4", "B"},
	}
	for i, r := range routes {
		assert.Equal(t, wantRoutesIds[i], nodeIds(r.ns))
	}

	a := findNode(G, "A")
	assert.ElementsMatch(t, []string{"B", "E", "F"}, outIds(a))

	g := findNode(G, "G")
	assert.ElementsMatch(t, []string{"A", "B", "C"}, outIds(g))

	b := findNode(G, "B")
	assert.ElementsMatch(t, []string{"A", "C", "G"}, inIds(b))

	e := findNode(G, "E")
	assert.ElementsMatch(t, []string{"A", "B"}, inIds(e))
}

func isEdge(e *graph.Edge, from, to string) bool {
	return e.From.ID == from && e.To.ID == to
}

func setLayer(g *graph.DGraph, n *graph.Node) {
	var l int
	switch n.ID {
	case "A", "C":
		l = 1
	case "B":
		l = 2
	case "D", "E":
		l = 3
	case "F":
		l = 4
	case "G":
		l = 0
	}
	if g.Layers[l] == nil {
		g.Layers[l] = &graph.Layer{}
	}
	g.Layers[l].Index = l
	g.Layers[l].Nodes = append(g.Layers[l].Nodes, n)
	n.Layer = l
}

func findNode(g *graph.DGraph, id string) *graph.Node {
	for _, n := range g.Nodes {
		if n.ID == id {
			return n
		}
	}
	return nil
}

func nodeIds(ns []*graph.Node) []string {
	ids := make([]string, len(ns))
	for i, n := range ns {
		ids[i] = n.ID
	}
	return ids
}

func outIds(n *graph.Node) []string {
	ids := make([]string, len(n.Out))
	for i, e := range n.Out {
		ids[i] = e.To.ID
	}
	return ids
}

func inIds(n *graph.Node) []string {
	ids := make([]string, len(n.In))
	for i, e := range n.In {
		ids[i] = e.From.ID
	}
	return ids
}

func fromEdgeSlice(es [][]string) *graph.DGraph {
	g := &graph.DGraph{}
	graph.EdgeSlice(es).Populate(g)
	return g
}

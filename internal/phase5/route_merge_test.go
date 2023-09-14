//go:build unit

package phase5

import (
	"sort"
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/internal/testfiles"
	"github.com/stretchr/testify/assert"
)

func TestMergeLongEdges(t *testing.T) {
	G := graph.FromEdgeSlice(testfiles.LongEdges)
	for _, e := range G.Edges {
		if isEdge(e, "C", "G") || isEdge(e, "B", "C") || isEdge(e, "B", "G") {
			e.Reverse()
		}
	}
	G.Layers = make(graph.Layers, 5)
	for _, n := range G.Nodes {
		setLayer(G, n)
	}

	assert.Len(t, G.Edges, 10)
	G.BreakLongEdges()
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

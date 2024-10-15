package ns

import (
	"testing"

	egraph "github.com/nulab/autog/graph"
	"github.com/nulab/autog/internal/graph"
	"github.com/stretchr/testify/assert"
)

func TestSpanningTree(t *testing.T) {
	g := &graph.DGraph{}
	egraph.EdgeSlice([][]string{
		{"a", "b"},
		{"b", "d"},
		{"b", "e"},
		{"a", "c"},
		{"c", "f"},
		{"c", "g"},
		{"c", "h"},
		{"f", "i"},
	}).Populate(g)

	for _, e := range g.Edges {
		e.IsInSpanningTree = true
	}
	p := Processor{
		lim: make(graph.NodeIntMap),
		low: make(graph.NodeIntMap),
	}
	p.setStreeValues(findNode(g, "a"))

	t.Run("low and lim", func(t *testing.T) {
		type tc struct {
			id       string
			low, lim int
		}
		tcs := []tc{
			{"a", 1, 9},
			{"b", 1, 3},
			{"c", 4, 8},
			{"d", 1, 1},
			{"e", 2, 2},
			{"f", 4, 5},
			{"g", 6, 6},
			{"h", 7, 7},
			{"i", 4, 4},
		}

		for _, tc := range tcs {
			n := findNode(g, tc.id)
			assert.Equalf(t, tc.low, p.low[n], "wrong low for n: %s", n.ID)
			assert.Equalf(t, tc.lim, p.lim[n], "wrong lim for n: %s", n.ID)
		}
	})

	t.Run("node in head component", func(t *testing.T) {
		e := findEdge(g, "c", "f")
		assert.True(t, p.inHeadComponent(findNode(g, "i"), e))
		assert.True(t, p.inHeadComponent(findNode(g, "f"), e))
		assert.False(t, p.inHeadComponent(findNode(g, "c"), e))
		assert.False(t, p.inHeadComponent(findNode(g, "d"), e))
		assert.False(t, p.inHeadComponent(findNode(g, "e"), e))
		assert.False(t, p.inHeadComponent(findNode(g, "h"), e))

		e.Reverse()
		assert.False(t, p.inHeadComponent(findNode(g, "i"), e))
		assert.False(t, p.inHeadComponent(findNode(g, "f"), e))
		assert.True(t, p.inHeadComponent(findNode(g, "c"), e))
		assert.True(t, p.inHeadComponent(findNode(g, "d"), e))
		assert.True(t, p.inHeadComponent(findNode(g, "e"), e))
		assert.True(t, p.inHeadComponent(findNode(g, "h"), e))

		e = findEdge(g, "a", "b")
		assert.False(t, p.inHeadComponent(findNode(g, "i"), e))
		assert.False(t, p.inHeadComponent(findNode(g, "f"), e))
		assert.False(t, p.inHeadComponent(findNode(g, "c"), e))
		assert.False(t, p.inHeadComponent(findNode(g, "a"), e))
		assert.True(t, p.inHeadComponent(findNode(g, "d"), e))
		assert.True(t, p.inHeadComponent(findNode(g, "e"), e))
		assert.False(t, p.inHeadComponent(findNode(g, "h"), e))

		e.Reverse()
		assert.True(t, p.inHeadComponent(findNode(g, "i"), e))
		assert.True(t, p.inHeadComponent(findNode(g, "f"), e))
		assert.True(t, p.inHeadComponent(findNode(g, "c"), e))
		assert.True(t, p.inHeadComponent(findNode(g, "a"), e))
		assert.False(t, p.inHeadComponent(findNode(g, "d"), e))
		assert.False(t, p.inHeadComponent(findNode(g, "e"), e))
		assert.True(t, p.inHeadComponent(findNode(g, "h"), e))

		e = findEdge(g, "b", "e")
		for _, n := range g.Nodes {
			if n.ID == "e" {
				assert.True(t, p.inHeadComponent(n, e))
			} else {
				assert.False(t, p.inHeadComponent(n, e))
			}
		}
	})
}

func findNode(g *graph.DGraph, id string) *graph.Node {
	for _, n := range g.Nodes {
		if n.ID == id {
			return n
		}
	}
	return nil
}

func findEdge(g *graph.DGraph, from, to string) *graph.Edge {
	for _, e := range g.Edges {
		if e.From.ID == from && e.To.ID == to {
			return e
		}
	}
	return nil
}

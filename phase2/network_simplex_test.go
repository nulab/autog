package phase2

import (
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/stretchr/testify/assert"
)

func TestPostorderTraversal(t *testing.T) {
	g := graph.FromEdgeSlice([][]string{
		{"a", "b"},
		{"b", "d"},
		{"b", "e"},
		{"a", "c"},
		{"c", "f"},
		{"c", "g"},
		{"c", "h"},
		{"f", "i"},
	})
	for _, e := range g.Edges {
		e.IsInSpanningTree = true
	}
	p := newNsProcessor()
	p.postOrderTraversal(findNode(g, "a"), graph.EdgeSet{}, 1)

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
}

func newNsProcessor() *networkSimplexProcessor {
	return &networkSimplexProcessor{
		lim: graph.NodeIntMap{},
		low: graph.NodeIntMap{},
	}
}

func findNode(g *graph.DGraph, id string) *graph.Node {
	for _, n := range g.Nodes {
		if n.ID == id {
			return n
		}
	}
	return nil
}

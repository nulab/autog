//go:build unit

package phase2

import (
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/internal/testfiles"
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
	p.setStreeValues(findNode(g, "a"))

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

func TestNSLayering(t *testing.T) {
	g := graph.FromEdgeSlice(testfiles.DotAbstract)
	execNetworkSimplex(g, graph.Params{NetworkSimplexThoroughness: 28, NetworkSimplexBalance: 1})

	want := expectedLayersAbstract()
	for _, n := range g.Nodes {
		if n.IsVirtual {
			continue
		}
		assert.Equalf(t, want[n.ID], n.Layer, "node %s layer %d but should be %d", n.ID, n.Layer, want[n.ID])
	}
}

func expectedLayersAbstract() map[string]int {
	// in dot the nodes 39 and 41 end up inverted
	// this is likely due to a different process order in the vbalance step
	// dot uses qsort which is unstable for equal values
	return map[string]int{
		"S1": 0, "S35": 0,
		"10": 1, "2": 1, "37": 1, "36": 1, "43": 1, "S24": 1,
		"S30": 2, "13": 2, "17": 2, "39": 4, "40": 2, "9": 2, "38": 2, "25": 2,
		"33": 3, "12": 3, "16": 3, "19": 3, "42": 3, "11": 3, "3": 3, "26": 3, "27": 3,
		"34": 4, "18": 4, "41": 2, "28": 4, "31": 4, "14": 4, "20": 4, "21": 4, "4": 4,
		"29": 5, "32": 5, "15": 5, "22": 5, "5": 5,
		"T30": 6, "23": 6, "T35": 6, "6": 6,
		"T1": 7, "T24": 7, "7": 7,
		"T8": 8,
	}
}

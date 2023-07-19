package cyclebreaking

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/autog/graph"
	"github.com/vibridi/autog/internal/testfiles"
)

func TestDepthFirst(t *testing.T) {
	testgs := testfiles.ReadTestDir("../internal/testfiles/elk_original")
	for _, g := range testgs {
		t.Run(g.Name, func(t *testing.T) {
			dg := graph.FromAdjacencyList(g.AdjacencyList())

			execDepthFirst(dg)

			assert.False(t, dg.HasCycles())
			// printReversedEdges(dg)
		})
	}
}

func printReversedEdges(g *graph.DGraph) {
	for _, n := range g.Nodes {
		for _, e := range n.Edges() {
			if e.IsReversed {
				e.Reverse()
				fmt.Println(e)
				e.Reverse()
			}
		}
	}
}

package cyclebreaking

import (
	"fmt"
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/internal/testfiles"
	"github.com/stretchr/testify/assert"
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

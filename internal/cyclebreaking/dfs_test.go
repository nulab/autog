package cyclebreaking

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/autog/internal/graph"
	"github.com/vibridi/autog/internal/testutils"
)

func TestDepthFirst(t *testing.T) {
	testgs := testutils.ReadTestDir("../testutils/elk/cyclic")
	for _, g := range testgs {
		t.Run(g.Name, func(t *testing.T) {
			dg := graph.FromAdjacencyList(g.AdjacencyList())

			DepthFirst.Process(dg)
			DepthFirst.Cleanup()

			assert.False(t, graph.HasCycles(dg))
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

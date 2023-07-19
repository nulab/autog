package cyclebreaking

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/autog/internal/graph"
	"github.com/vibridi/autog/internal/testutils"
)

func TestGreedy(t *testing.T) {
	testgs := testutils.ReadTestDir("../testutils/elk/cyclic")
	for _, g := range testgs {
		t.Run(g.Name, func(t *testing.T) {
			dg := graph.FromAdjacencyList(g.AdjacencyList())

			Greedy.Process(dg)
			Greedy.Cleanup()

			assert.False(t, graph.HasCycles(dg))
			printReversedEdges(dg)
		})
	}
}

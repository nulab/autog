package cyclebreaking

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/autog/graph"
	"github.com/vibridi/autog/internal/testutils"
)

func TestGreedy(t *testing.T) {
	testgs := testutils.ReadTestDir("../internal/testutils/elk/cyclic")
	for _, g := range testgs {
		t.Run(g.Name, func(t *testing.T) {
			dg := graph.FromAdjacencyList(g.AdjacencyList())

			execGreedy(dg)

			assert.False(t, graph.HasCycles(dg))
			printReversedEdges(dg)
		})
	}
}

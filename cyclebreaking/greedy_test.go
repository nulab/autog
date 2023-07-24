package cyclebreaking

import (
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/internal/testfiles"
	"github.com/stretchr/testify/assert"
)

func TestGreedy(t *testing.T) {
	testgs := testfiles.ReadTestDir("../internal/testfiles/elk_original")
	for _, g := range testgs {
		t.Run(g.Name, func(t *testing.T) {
			dg := graph.FromAdjacencyList(g.AdjacencyList())

			execGreedy(dg)

			assert.False(t, dg.HasCycles())
			printReversedEdges(dg)
		})
	}
}

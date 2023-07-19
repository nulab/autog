package cyclebreaking

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/autog/graph"
	"github.com/vibridi/autog/internal/testfiles"
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

package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/autog/internal/testfiles"
)

func TestHasCycles(t *testing.T) {
	testgs := testfiles.ReadTestDir("../internal/testfiles/elk_original")
	for _, g := range testgs {
		t.Run(g.Name, func(t *testing.T) {
			dg := FromAdjacencyList(g.AdjacencyList())
			assert.True(t, dg.HasCycles())
		})
	}
}

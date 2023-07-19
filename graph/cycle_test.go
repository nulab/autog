package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/autog/internal/testutils"
)

func TestHasCycles(t *testing.T) {
	testgs := testutils.ReadTestDir("../testutils/elk/cyclic")
	for _, g := range testgs {
		t.Run(g.Name, func(t *testing.T) {
			dg := FromAdjacencyList(g.AdjacencyList())
			assert.True(t, HasCycles(dg))
		})
	}
}

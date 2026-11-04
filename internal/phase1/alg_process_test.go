package phase1

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nulab/autog/internal/graph"
)

func TestRemoveTwoNodeCycles(t *testing.T) {
	t.Run("does not reverse parallel edges", func(t *testing.T) {
		es := graph.EdgeSlice{{"a", "b"}, {"a", "b"}}
		g := &graph.DGraph{}
		es.Populate(g)
		removeTwoNodeCycles(g)
		assert.False(t, g.Edges[0].IsReversed)
		assert.False(t, g.Edges[1].IsReversed)
	})

	t.Run("reverses actual two node cycles", func(t *testing.T) {
		es := graph.EdgeSlice{{"a", "b"}, {"b", "c"}, {"a", "b"}, {"c", "b"}}
		g := &graph.DGraph{}
		es.Populate(g)
		removeTwoNodeCycles(g)
		assert.False(t, g.Edges[0].IsReversed)
		assert.False(t, g.Edges[1].IsReversed)
		assert.False(t, g.Edges[2].IsReversed)
		assert.True(t, g.Edges[3].IsReversed)
	})
}

//go:build unit

package testfiles

import (
	"testing"

	"github.com/nulab/autog"
	"github.com/nulab/autog/graph"
	ig "github.com/nulab/autog/internal/graph"
	"github.com/stretchr/testify/assert"
)

func TestOutputVirtualNodes(t *testing.T) {
	t.Run("keep virtual nodes", func(t *testing.T) {
		g := &ig.DGraph{}
		graph.EdgeSlice(simpleVirtualNodes).Populate(g)
		layout := autog.Layout(
			g,
			autog.WithPositioning(autog.PositioningVAlign),
			autog.WithEdgeRouting(autog.EdgeRoutingNoop),
			autog.WithOutputVirtualNodes(true),
		)
		assert.Len(t, layout.Nodes, 4)
		assert.Len(t, layout.Edges, 3)
	})

	t.Run("clip output nodes", func(t *testing.T) {
		g := &ig.DGraph{}
		graph.EdgeSlice(simpleVirtualNodes).Populate(g)
		layout := autog.Layout(
			g,
			autog.WithPositioning(autog.PositioningVAlign),
			autog.WithEdgeRouting(autog.EdgeRoutingNoop),
			autog.WithOutputVirtualNodes(false),
		)
		assert.Equal(t, 3, len(layout.Nodes))
		assert.Equal(t, 3, cap(layout.Nodes))
		assert.Equal(t, 4, len(g.Nodes))

		assert.Len(t, layout.Edges, 3)
	})
}

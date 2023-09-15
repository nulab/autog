package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdge(t *testing.T) {
	n1 := &Node{ID: "1", Layer: 0, LayerPos: 0}
	n2 := &Node{ID: "2", Layer: 1, LayerPos: 1}
	e := NewEdge(n1, n2, 0)

	n1.Out = []*Edge{e}
	n2.In = []*Edge{e}

	// normal direction
	assert.Equal(t, 0, n1.Indeg())
	assert.Equal(t, 1, n1.Outdeg())

	assert.Equal(t, 1, n2.Indeg())
	assert.Equal(t, 0, n2.Outdeg())

	// reverse
	e.Reverse()
	assert.Equal(t, 1, n1.Indeg())
	assert.Equal(t, 0, n1.Outdeg())

	assert.Equal(t, 0, n2.Indeg())
	assert.Equal(t, 1, n2.Outdeg())

	// reverse it back
	e.Reverse()
	assert.Equal(t, 0, n1.Indeg())
	assert.Equal(t, 1, n1.Outdeg())

	assert.Equal(t, 1, n2.Indeg())
	assert.Equal(t, 0, n2.Outdeg())

	n3 := &Node{ID: "3", Layer: 0, LayerPos: 1}
	n4 := &Node{ID: "4", Layer: 1, LayerPos: 0}
	f := NewEdge(n3, n4, 0)

	n3.Out = []*Edge{f}
	n4.In = []*Edge{f}

	assert.True(t, e.Crosses(f))
	e.Reverse()
	assert.True(t, e.Crosses(f))
}

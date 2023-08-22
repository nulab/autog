package phase3

import (
	"fmt"
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/stretchr/testify/assert"
)

var nodes = map[int]*graph.Node{
	0: {ID: "n0", Layer: 0},
	1: {ID: "n1", Layer: 1},
	2: {ID: "n2", Layer: 1},
	3: {ID: "n3", Layer: 1},
	4: {ID: "n4", Layer: 1},
	5: {ID: "n5", Layer: 1},
	6: {ID: "n6", Layer: 2},
	7: {ID: "n7", Layer: 2},
	8: {ID: "n8", Layer: 3},
}

func TestInitFlatEdges(t *testing.T) {
	t.Run("ordered sequence", func(t *testing.T) {
		edges := []*graph.Edge{
			graph.NewEdge(nodes[1], nodes[2], 0),
			graph.NewEdge(nodes[3], nodes[4], 0),
			graph.NewEdge(nodes[4], nodes[5], 0),
			graph.NewEdge(nodes[7], nodes[8], 0),
		}
		fpos := initFixedPositions(edges)
		assert.Len(t, fpos.mustAfter, 3)
		assert.Len(t, fpos.mustBefore, 3)
		assert.Equal(t, nodes[1], fpos.mustAfter[nodes[2]])
		assert.Equal(t, nodes[4], fpos.mustAfter[nodes[5]])
		assert.Equal(t, nodes[3], fpos.mustAfter[nodes[4]])

		n, i := fpos.head(nodes[5])
		assert.Equal(t, nodes[3], n)
		assert.Equal(t, 2, i)
	})

	t.Run("unordered sequence", func(t *testing.T) {
		edges := []*graph.Edge{
			graph.NewEdge(nodes[4], nodes[5], 0),
			graph.NewEdge(nodes[1], nodes[2], 0),
			graph.NewEdge(nodes[7], nodes[8], 0),
			graph.NewEdge(nodes[3], nodes[4], 0),
		}
		fpos := initFixedPositions(edges)
		assert.Len(t, fpos.mustAfter, 3)
		assert.Len(t, fpos.mustBefore, 3)
		assert.Equal(t, nodes[2], fpos.mustBefore[nodes[1]])
		assert.Equal(t, nodes[5], fpos.mustBefore[nodes[4]])
		assert.Equal(t, nodes[4], fpos.mustBefore[nodes[3]])

		n, i := fpos.head(nodes[5])
		assert.Equal(t, nodes[3], n)
		assert.Equal(t, 2, i)
	})

	t.Run("same source", func(t *testing.T) {
		edges := []*graph.Edge{
			graph.NewEdge(nodes[1], nodes[3], 0),
			graph.NewEdge(nodes[7], nodes[8], 0),
			graph.NewEdge(nodes[3], nodes[4], 0),
			graph.NewEdge(nodes[3], nodes[5], 0),
		}
		fpos := initFixedPositions(edges)
		assert.Len(t, fpos.mustAfter, 3)
		assert.Len(t, fpos.mustBefore, 3)
		assert.Equal(t, nodes[3], fpos.mustBefore[nodes[1]])
		assert.Equal(t, nodes[4], fpos.mustBefore[nodes[3]])
		assert.Equal(t, nodes[5], fpos.mustBefore[nodes[4]])

		fmt.Println(fpos.mustAfter)
		fmt.Println(fpos.mustBefore)

		n, i := fpos.head(nodes[5])
		assert.Equal(t, nodes[1], n)
		assert.Equal(t, 3, i)
	})

	t.Run("same target", func(t *testing.T) {
		edges := []*graph.Edge{
			graph.NewEdge(nodes[1], nodes[3], 0),
			graph.NewEdge(nodes[7], nodes[8], 0),
			graph.NewEdge(nodes[3], nodes[5], 0),
			graph.NewEdge(nodes[4], nodes[5], 0),
		}
		fpos := initFixedPositions(edges)
		fmt.Println(fpos.mustAfter)
		fmt.Println(fpos.mustBefore)

		assert.Len(t, fpos.mustAfter, 3)
		assert.Len(t, fpos.mustBefore, 3)
		assert.Equal(t, nodes[3], fpos.mustBefore[nodes[1]])
		assert.Equal(t, nodes[4], fpos.mustBefore[nodes[3]])
		assert.Equal(t, nodes[5], fpos.mustBefore[nodes[4]])

		n, i := fpos.head(nodes[5])
		assert.Equal(t, nodes[1], n)
		assert.Equal(t, 3, i)
	})

	t.Run("success", func(t *testing.T) {
		edges := []*graph.Edge{
			graph.NewEdge(nodes[1], nodes[3], 0),
			graph.NewEdge(nodes[6], nodes[7], 0),
			graph.NewEdge(nodes[3], nodes[5], 0),
			graph.NewEdge(nodes[4], nodes[5], 0),
			graph.NewEdge(nodes[1], nodes[4], 0), // should be no op
		}
		fpos := initFixedPositions(edges)
		fmt.Println(fpos.mustAfter)
		fmt.Println(fpos.mustBefore)

		assert.Len(t, fpos.mustAfter, 4)
		assert.Len(t, fpos.mustBefore, 4)
		assert.Equal(t, nodes[3], fpos.mustBefore[nodes[1]])
		assert.Equal(t, nodes[4], fpos.mustBefore[nodes[3]])
		assert.Equal(t, nodes[5], fpos.mustBefore[nodes[4]])

		assert.Empty(t, fpos.mustAfter[nodes[1]])
		assert.Equal(t, nodes[1], fpos.mustAfter[nodes[3]])
		assert.Equal(t, nodes[3], fpos.mustAfter[nodes[4]])
		assert.Equal(t, nodes[4], fpos.mustAfter[nodes[5]])

		n, i := fpos.head(nodes[5])
		assert.Equal(t, nodes[1], n)
		assert.Equal(t, 3, i)
	})
}

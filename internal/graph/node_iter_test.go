package graph

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisitEdges(t *testing.T) {
	n := &Node{ID: "N"}
	n.In = []*Edge{
		{edge: edge{From: &Node{ID: "A1"}, To: n}},
		{edge: edge{From: &Node{ID: "A2"}, To: n}},
		{edge: edge{From: &Node{ID: "A3"}, To: n}},
	}
	n.Out = []*Edge{
		{edge: edge{From: n, To: &Node{ID: "B1"}}},
		{edge: edge{From: n, To: &Node{ID: "B2"}}},
		{edge: edge{From: n, To: &Node{ID: "B3"}}},
		{edge: edge{From: n, To: &Node{ID: "B4"}}},
	}

	i := 0
	n.VisitEdges(func(e *Edge) {
		if i < 3 {
			assert.True(t, strings.HasPrefix(e.From.ID, "A"))
		} else {
			assert.True(t, strings.HasPrefix(e.To.ID, "B"))
		}
		i++
	})
	assert.Equal(t, 7, i)
}

package graph

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisitEdges(t *testing.T) {
	ids := strings.Split("abcde", "")
	es := []*Edge{}
	for _, id := range ids {
		es = append(es, &Edge{edge: edge{ID: id}})
	}

	n := &Node{
		In:  es[:3],
		Out: es[3:],
	}

	i := 0
	n.VisitEdges(func(e *Edge) {
		assert.Equal(t, ids[i], e.ID)
		i++
	})
	assert.Equal(t, 5, i)
}

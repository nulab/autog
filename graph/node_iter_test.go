package graph

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode_EdgeIter(t *testing.T) {
	ids := strings.Split("abcde", "")
	es := []*Edge{}
	for _, id := range ids {
		es = append(es, &Edge{edge: edge{ID: id}})
	}

	n := &Node{
		In:  es[:3],
		Out: es[3:],
	}

	itr := n.EdgeIter()
	for i := 0; itr.HasNext(); i++ {
		e := itr.Next()
		assert.Equal(t, ids[i], e.ID)
	}
	assert.False(t, itr.HasNext())
	assert.Equal(t, 5, itr.i)
}

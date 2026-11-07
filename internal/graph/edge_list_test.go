package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdgeListRemove(t *testing.T) {
	a := &Edge{}
	b := &Edge{}
	c := &Edge{}

	t.Run("single match", func(t *testing.T) {
		list := EdgeList{a, b, c}
		list.Remove(b)
		assert.Equal(t, EdgeList{a, c}, list)
	})

	t.Run("adjacent matches", func(t *testing.T) {
		list := EdgeList{a, a, b}
		list.Remove(a)
		assert.Equal(t, EdgeList{b}, list)
	})

	t.Run("non-adjacent matches", func(t *testing.T) {
		list := EdgeList{a, b, a, c, a}
		list.Remove(a)
		assert.Equal(t, EdgeList{b, c}, list)
	})

	t.Run("no match", func(t *testing.T) {
		list := EdgeList{a, b}
		list.Remove(c)
		assert.Equal(t, EdgeList{a, b}, list)
	})

	t.Run("nil list", func(t *testing.T) {
		var list *EdgeList
		assert.NotPanics(t, func() { list.Remove(a) })
	})
}

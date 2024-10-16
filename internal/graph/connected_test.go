package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponents(t *testing.T) {
	t.Run("one component", func(t *testing.T) {
		es := [][]string{
			{"a", "b"},
			{"b", "c"},
		}
		g := &DGraph{}
		EdgeSlice(es).Populate(g)

		comp := g.ConnectedComponents()
		assert.Len(t, comp, 1)
		assert.True(t, comp[0] == g)
	})

	t.Run("multiple components", func(t *testing.T) {
		es := [][]string{
			{"a", "b"},
			{"b", "c"},
			{"f", "g"},
		}
		g := &DGraph{}
		EdgeSlice(es).Populate(g)

		comp := g.ConnectedComponents()
		assert.Len(t, comp, 2)
		assert.ElementsMatch(t, []string{"a", "b", "c"}, ids(comp[0].Nodes))
		assert.ElementsMatch(t, []string{"f", "g"}, ids(comp[1].Nodes))
	})

	t.Run("self-loop", func(t *testing.T) {
		es := [][]string{
			{"a", "b"}, {"b", "c"},
			{"f", "g"}, {"g", "h"}, {"h", "i"},
			{"u", "u"},
			{"j", "k"},
			{"l", "j"},
			{"z", "f"},
		}
		g := &DGraph{}
		EdgeSlice(es).Populate(g)

		comp := g.ConnectedComponents()
		assert.Len(t, comp, 4)
		assert.ElementsMatch(t, []string{"a", "b", "c"}, ids(comp[0].Nodes))
		assert.ElementsMatch(t, []string{"f", "g", "h", "i", "z"}, ids(comp[1].Nodes))
		assert.ElementsMatch(t, []string{"u"}, ids(comp[2].Nodes))
		assert.ElementsMatch(t, []string{"j", "k", "l"}, ids(comp[3].Nodes))
	})

}

func ids(ns []*Node) (ids []string) {
	for _, n := range ns {
		ids = append(ids, n.ID)
	}
	return
}

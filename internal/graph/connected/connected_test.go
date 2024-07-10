package connected

import (
	"testing"

	"github.com/nulab/autog/graph"
	ig "github.com/nulab/autog/internal/graph"
	"github.com/stretchr/testify/assert"
)

func TestComponents(t *testing.T) {
	t.Run("one component", func(t *testing.T) {
		es := [][]string{
			{"a", "b"},
			{"b", "c"},
		}
		g := &ig.DGraph{}
		graph.EdgeSlice(es).Populate(g)

		comp := Components(g)
		assert.Len(t, comp, 1)
		assert.True(t, comp[0] == g)
	})

	t.Run("multiple components", func(t *testing.T) {
		es := [][]string{
			{"a", "b"},
			{"b", "c"},
			{"f", "g"},
		}
		g := &ig.DGraph{}
		graph.EdgeSlice(es).Populate(g)

		comp := Components(g)
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
		g := &ig.DGraph{}
		graph.EdgeSlice(es).Populate(g)

		comp := Components(g)
		assert.Len(t, comp, 4)
		assert.ElementsMatch(t, []string{"a", "b", "c"}, ids(comp[0].Nodes))
		assert.ElementsMatch(t, []string{"f", "g", "h", "i", "z"}, ids(comp[1].Nodes))
		assert.ElementsMatch(t, []string{"u"}, ids(comp[2].Nodes))
		assert.ElementsMatch(t, []string{"j", "k", "l"}, ids(comp[3].Nodes))
	})

}

func ids(ns []*ig.Node) (ids []string) {
	for _, n := range ns {
		ids = append(ids, n.ID)
	}
	return
}

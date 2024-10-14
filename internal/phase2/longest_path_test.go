package phase2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLongestPath(t *testing.T) {
	g := fromEdgeSlice([][]string{
		{"F", "B"},
		{"F", "Z"},
		{"B", "A"},
		{"B", "D"},
		{"B", "K"},
		{"Z", "a1"},
		{"Z", "b1"},
		{"A", "C"},
		{"C", "U"},
	})

	want := map[string]int{
		"F": 0,
		"B": 1, "Z": 1,
		"A": 2, "D": 2, "K": 2, "a1": 2, "b1": 2,
		"C": 3,
		"U": 4,
	}

	execLongestPath(g)
	for _, n := range g.Nodes {
		assert.Equal(t, want[n.ID], n.Layer)
	}
}

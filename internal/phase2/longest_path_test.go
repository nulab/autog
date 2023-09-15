package phase2

import (
	"fmt"
	"testing"

	"github.com/nulab/autog/internal/graph"
)

func TestLongestPath(t *testing.T) {
	g := graph.FromEdgeSlice([][]string{
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

	execLongestPath(g)
	for _, n := range g.Nodes {
		fmt.Println(n.ID, n.Layer)
	}
}

package graph

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/autog/internal/testutils"
)

func TestSpanningTree(t *testing.T) {
	g := testutils.ReadTestFile("../testutils/elk_constructed", "simple_acyclic.json")
	dg := FromAdjacencyList(g.AdjacencyList())
	assert.False(t, HasCycles(dg))

	s := dg.SpanningTree()
	fmt.Println(s)
}

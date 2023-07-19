package layering

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/autog/internal/cyclebreaking"
	"github.com/vibridi/autog/internal/graph"
	"github.com/vibridi/autog/internal/testutils"
)

func TestA(t *testing.T) {
	g := testutils.ReadTestFile("../testutils/elk_constructed", "simple_acyclic.json")
	dg := graph.FromAdjacencyList(g.AdjacencyList())
	assert.False(t, graph.HasCycles(dg))

	s := dg.SpanningTree()
	fmt.Println(s)
	ns := &networkSimplexProcessor{}
	ns.initLayers(dg)
	for _, n := range dg.Nodes {
		fmt.Println(n.ID, n.Layer)
	}
}

func TestB(t *testing.T) {
	g := testutils.ReadTestFile("../testutils/elk_constructed", "simple_acyclic.json")
	dg := graph.FromAdjacencyList(g.AdjacencyList())
	assert.False(t, graph.HasCycles(dg))

	ns := &networkSimplexProcessor{}
	ns.Process(dg)
	for _, n := range dg.Nodes {
		fmt.Println(n.ID, n.Layer)
	}
}

func TestLayering(t *testing.T) {
	testgs := testutils.ReadTestDir("../testutils/elk/cyclic")
	for _, g := range testgs {
		dg := graph.FromAdjacencyList(g.AdjacencyList())
		if graph.HasCycles(dg) {
			cyclebreaking.DepthFirst.Process(dg)
			cyclebreaking.DepthFirst.Cleanup()
		}
		t.Run(g.Name, func(t *testing.T) {
			NetworkSimplex.Process(dg)
			NetworkSimplex.Cleanup()

			for _, e := range dg.Edges {
				assert.True(t, e.From.Layer < e.To.Layer)
			}
		})
	}
}

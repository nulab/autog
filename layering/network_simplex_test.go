package layering

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/autog/cyclebreaking"
	"github.com/vibridi/autog/graph"
	"github.com/vibridi/autog/internal/testfiles"
)

func TestA(t *testing.T) {
	g := testfiles.ReadTestFile("../testfiles/elk_constructed", "simple_acyclic.json")
	dg := graph.FromAdjacencyList(g.AdjacencyList())
	assert.False(t, dg.HasCycles())

	s := dg.SpanningTree()
	fmt.Println(s)
	ns := &networkSimplexProcessor{}
	ns.initLayers(dg)
	for _, n := range dg.Nodes {
		fmt.Println(n.ID, n.Layer)
	}
}

func TestB(t *testing.T) {
	g := testfiles.ReadTestFile("../testfiles/elk_constructed", "simple_acyclic.json")
	dg := graph.FromAdjacencyList(g.AdjacencyList())
	assert.False(t, dg.HasCycles())

	execNetworkSimplex(dg)
	for _, n := range dg.Nodes {
		fmt.Println(n.ID, n.Layer)
	}
}

func TestLayering(t *testing.T) {
	testgs := testfiles.ReadTestDir("../internal/testfiles/elk_original")
	for _, g := range testgs {
		dg := graph.FromAdjacencyList(g.AdjacencyList())
		if dg.HasCycles() {
			cyclebreaking.DEPTH_FIRST.Process(dg)
		}
		t.Run(g.Name, func(t *testing.T) {
			for _, subg := range dg.ConnectedComponents() {
				execNetworkSimplex(subg)
				for _, e := range subg.Edges {
					assert.True(t, e.From.Layer < e.To.Layer)
				}
			}
		})
	}
}

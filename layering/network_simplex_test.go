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

func TestPostorderTraversal(t *testing.T) {
	p := &networkSimplexProcessor{
		poIndex: 1,
		lim:     make(graph.NodeMap),
		low:     make(graph.NodeMap),
	}
	g := testfiles.ReadTestFile("../internal/testfiles/elk_constructed", "simple_acyclic.json")
	dg := graph.FromAdjacencyList(g.AdjacencyList())
	dg.SpanningTree()
	fmt.Println(dg.Edges)
	fmt.Println("root", dg.Nodes[0])
	p.postOrderTraversal(dg.Nodes[0], graph.EdgeSet{})
	fmt.Println(p.lim)
	fmt.Println(p.low)
}

func TestNetworkSimplexLayering(t *testing.T) {
	testgs := testfiles.ReadTestDir("../internal/testfiles/elk_original")
	for _, g := range testgs {
		dg := graph.FromAdjacencyList(g.AdjacencyList())
		if dg.HasCycles() {
			cyclebreaking.DepthFirst.Process(dg)
		}
		t.Run(g.Name, func(t *testing.T) {
			for _, subg := range dg.ConnectedComponents() {
				t.Run("component:"+subg.Nodes[0].ID, func(t *testing.T) {
					execNetworkSimplex(subg)
					for _, e := range subg.Edges {
						assert.True(t, e.From.Layer < e.To.Layer)
					}
				})

			}
		})
	}
}

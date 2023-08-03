package phase2

import (
	"fmt"
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/internal/testfiles"
	"github.com/nulab/autog/phase1"
	"github.com/stretchr/testify/assert"
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
	testgs := testfiles.ReadTestDir("../internal/testfiles/elk_relabeled")
	for _, g := range testgs {
		if g.Name != "ci_router_ComplexRouter.json" {
			continue
		}
		dg := graph.FromAdjacencyList(g.AdjacencyList())
		if dg.HasCycles() {
			phase1.DepthFirst.Process(dg, nil)
		}
		t.Run(g.Name, func(t *testing.T) {
			for _, subg := range dg.ConnectedComponents() {
				t.Run("component:"+subg.Nodes[0].ID, func(t *testing.T) {
					execNetworkSimplex(subg)
					for _, e := range subg.Edges {
						assert.True(t, e.From.Layer <= e.To.Layer)
					}
				})

			}
		})
	}
}

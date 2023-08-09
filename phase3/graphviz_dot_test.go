package phase3

import (
	"fmt"
	"slices"
	"testing"

	"github.com/nulab/autog/graph"
)

// func TestVirtualNodes(t *testing.T) {
// 	g := graph.FromAdjacencyList(map[string][]string{
// 		"N1": {"N2", "N4"},
// 		"N2": {"N3", "N4"},
// 		"N3": {"N4"},
// 	})
// 	assignLayersVertical(g)
// 	g.BreakLongEdges()
//
// 	assert.ElementsMatch(t, []string{"N1", "N2", "N3", "N4", "V1", "V2", "V3"}, nodeIds(g))
// 	assert.ElementsMatch(t, []string{"N1->N2", "N2->N3", "N3->N4", "N1->V1", "V1->V3", "V3->N4", "N2->V2", "V2->N4"}, edgeStrings(g))
// }
//
// func TestCrossings(t *testing.T) {
// 	g := graph.FromAdjacencyList(map[string][]string{
// 		"N1": {"N6"},
// 		"N2": {"N5"},
// 		"N3": {"N6", "N8"},
// 		"N4": {"N6", "N7"},
// 	})
// 	for _, n := range g.Nodes {
// 		switch n.ID {
// 		case "N1", "N2", "N3", "N4":
// 			n.Layer = 0
// 		default:
// 			n.Layer = 1
// 		}
// 	}
// 	nodes := g.Nodes
// 	sort.Slice(nodes, func(i, j int) bool {
// 		return nodes[i].ID < nodes[j].ID
// 	})
//
// 	layers := map[int]*graph.Layer{
// 		0: {Index: 0, Nodes: []*graph.Node{nodes[0], nodes[1], nodes[2], nodes[3]}},
// 		1: {Index: 1, Nodes: []*graph.Node{nodes[4], nodes[5], nodes[6], nodes[7]}},
// 	}
//
// 	p := &graphvizDotProcessor{positions: graph.NodeIntMap{}}
//
// 	for _, l := range layers {
// 		for i, n := range l.Nodes {
// 			p.setPos(n, i)
// 		}
// 	}
//
// 	assert.Equal(t, 3, crossings(layers))
// 	p.swap(nodes[6], nodes[7])
// 	assert.Equal(t, 2, crossings(layers))
// 	p.swap(nodes[6], nodes[7])
// 	p.swap(nodes[4], nodes[5])
// 	assert.Equal(t, 4, crossings(layers))
// 	p.swap(nodes[4], nodes[6])
// 	assert.Equal(t, 5, crossings(layers))
// 	p.swap(nodes[4], nodes[7])
// 	assert.Equal(t, 6, crossings(layers))
// }

func TestAAA(t *testing.T) {
	g := buildGraph()
	execGraphvizDot(g, graph.Params{})
}

// func TestGansnerNorthOrdering(t *testing.T) {
// 	testgs := testfiles.ReadTestDir("../internal/testfiles/elk_relabeled")
// 	for _, g := range testgs {
// 		dg := graph.FromElk(g)
// 		if dg.HasCycles() {
// 			phase1.DepthFirst.Process(dg, graph.Params{})
// 		}
// 		t.Run(g.Name, func(t *testing.T) {
// 			if len(g.Nodes) >= 100 {
// 				t.Skip()
// 			}
// 			for _, subg := range dg.ConnectedComponents() {
// 				t.Run("component:"+subg.Nodes[0].ID, func(t *testing.T) {
// 					phase2.NetworkSimplex.Process(subg, graph.Params{})
// 					execGraphvizDot(subg, graph.Params{})
//
// 					indices := map[int]map[int]bool{}
// 					for _, n := range subg.Nodes {
// 						m := indices[n.Layer]
// 						if m == nil {
// 							m = map[int]bool{}
// 							indices[n.Layer] = m
// 						}
// 						assert.False(t, m[n.LayerPos])
// 						m[n.LayerPos] = true
// 					}
// 					// printNodeOrders(subg)
// 				})
//
// 			}
// 		})
// 	}
// }

func assignLayersVertical(g *graph.DGraph) {
	for _, n := range g.Nodes {
		f := func(id string) int {
			switch id {
			case "N1":
				return 0
			case "N2":
				return 1
			case "N3":
				return 2
			case "N4":
				return 3
			}
			return -1
		}
		n.Layer = f(n.ID)
	}
}

func nodeIds(g *graph.DGraph) []string {
	ids := make([]string, len(g.Nodes))
	for i, n := range g.Nodes {
		ids[i] = n.ID
	}
	return ids
}

func edgeStrings(g *graph.DGraph) []string {
	ids := make([]string, len(g.Edges))
	for i, e := range g.Edges {
		ids[i] = fmt.Sprintf("%s->%s", e.From, e.To)
	}
	return ids
}

func printNodeOrders(g *graph.DGraph) {
	slices.SortFunc(g.Nodes, func(a, b *graph.Node) int {
		if a.Layer != b.Layer {
			return a.Layer - b.Layer
		}
		return a.LayerPos - b.LayerPos
	})
	for _, n := range g.Nodes {
		fmt.Printf("%s L:%d I:%d\n", n.ID, n.Layer, n.LayerPos)
	}
}

func buildGraph() *graph.DGraph {
	type edgelist = []*graph.Edge

	n10 := &graph.Node{ID: "N10"}
	n9 := &graph.Node{ID: "N9"}
	n15 := &graph.Node{ID: "N15"}
	n8 := &graph.Node{ID: "N8"}
	// n3 := &graph.Node{ID: "N3"}
	n6 := &graph.Node{ID: "N6"}
	n1 := &graph.Node{ID: "N1"}
	n4 := &graph.Node{ID: "N4"}
	n2 := &graph.Node{ID: "N2"}
	n16 := &graph.Node{ID: "N16"}
	n12 := &graph.Node{ID: "N12"}
	n13 := &graph.Node{ID: "N13"}
	n14 := &graph.Node{ID: "N14"}
	n5 := &graph.Node{ID: "N5"}
	n7 := &graph.Node{ID: "N7"}
	n11 := &graph.Node{ID: "N11"}
	v1 := &graph.Node{ID: "V1"}
	v2 := &graph.Node{ID: "V2"}

	// N10
	n9n10 := graph.NewEdge(n9, n10, 0)
	n15n10 := graph.NewEdge(n15, n10, 0)
	n10.In = edgelist{n9n10, n15n10}

	n10n11 := graph.NewEdge(n10, n11, 0)
	n10.Out = edgelist{n10n11}

	// N9
	n15n9 := graph.NewEdge(n15, n9, 0)
	n9.In = edgelist{n15n9}
	n9.Out = edgelist{n9n10}

	// N15
	n8n15 := graph.NewEdge(n8, n15, 0)
	n15n1 := graph.NewEdge(n15, n1, 0)
	n15.In = edgelist{n8n15}
	n15.Out = edgelist{n15n1, n15n9, n15n10}

	// N8
	// n3n8 := graph.NewEdge(n3, n8, 0)

	g := &graph.DGraph{
		Nodes: []*graph.Node{n10, n9, n15, n8, n6, n1, n4, n2, n16, n12, n13, n14, n5, n7, n11, v1, v2},
	}
	return g
}

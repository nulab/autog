package ordering

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/autog/cyclebreaking"
	"github.com/vibridi/autog/graph"
	"github.com/vibridi/autog/internal/testfiles"
	"github.com/vibridi/autog/layering"
	"golang.org/x/exp/slices"
)

func TestVirtualNodes(t *testing.T) {
	g := graph.FromAdjacencyList(map[string][]string{
		"N1": {"N2", "N4"},
		"N2": {"N3", "N4"},
		"N3": {"N4"},
	})
	assignLayersVertical(g)
	breakLongEdges(g)

	assert.ElementsMatch(t, []string{"N1", "N2", "N3", "N4", "V1", "V2", "V3"}, nodeIds(g))
	assert.ElementsMatch(t, []string{"N1->N2", "N2->N3", "N3->N4", "N1->V1", "V1->V3", "V3->N4", "N2->V2", "V2->N4"}, edgeStrings(g))
}

func TestCrossings(t *testing.T) {
	g := graph.FromAdjacencyList(map[string][]string{
		"N1": {"N6"},
		"N2": {"N5"},
		"N3": {"N6", "N8"},
		"N4": {"N6", "N7"},
	})
	for _, n := range g.Nodes {
		switch n.ID {
		case "N1", "N2", "N3", "N4":
			n.Layer = 0
		default:
			n.Layer = 1
		}
	}
	nodes := g.Nodes
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ID < nodes[j].ID
	})

	p := &gansnerNorthProcessor{
		layers: map[int][]*graph.Node{
			0: {nodes[0], nodes[1], nodes[2], nodes[3]},
			1: {nodes[4], nodes[5], nodes[6], nodes[7]},
		},
		orders: map[*graph.Node]int{
			nodes[0]: 0, nodes[1]: 1, nodes[2]: 2, nodes[3]: 3,
			nodes[4]: 0, nodes[5]: 1, nodes[6]: 2, nodes[7]: 3,
		},
		minL: 0,
		maxL: 1,
	}

	assert.Equal(t, 2, p.crossingsOf(1, nodes[6], nodes[7]))
	assert.Equal(t, 1, p.crossingsOf(1, nodes[7], nodes[6]))

	assert.Equal(t, 3, p.crossingsOf(1, nodes[5], nodes[6]))
	assert.Equal(t, 4, p.crossingsOf(1, nodes[6], nodes[5]))

	assert.Equal(t, 2, p.crossingsOf(1, nodes[4], nodes[5]))
	assert.Equal(t, 3, p.crossingsOf(1, nodes[5], nodes[4]))

	assert.Equal(t, 2, p.crossingsOf(0, nodes[2], nodes[3]))
	assert.Equal(t, 1, p.crossingsOf(0, nodes[3], nodes[2]))

	assert.Equal(t, 1, p.crossingsOf(0, nodes[0], nodes[1]))
	assert.Equal(t, 0, p.crossingsOf(0, nodes[1], nodes[0]))
}

func TestGansnerNorthOrdering(t *testing.T) {
	testgs := testfiles.ReadTestDir("../internal/testfiles/elk_relabeled")
	for _, g := range testgs {
		dg := graph.FromAdjacencyList(g.AdjacencyList())
		if dg.HasCycles() {
			cyclebreaking.DepthFirst.Process(dg)
		}
		t.Run(g.Name, func(t *testing.T) {
			if len(g.Nodes) >= 100 {
				t.Skip()
			}
			for _, subg := range dg.ConnectedComponents() {
				t.Run("component:"+subg.Nodes[0].ID, func(t *testing.T) {
					layering.NetworkSimplex.Process(subg)
					execGansnerNorth(subg)

					indices := map[int]map[int]bool{}
					for _, n := range subg.Nodes {
						m := indices[n.Layer]
						if m == nil {
							m = map[int]bool{}
							indices[n.Layer] = m
						}
						assert.False(t, m[n.LayerIdx])
						m[n.LayerIdx] = true
					}
					// printNodeOrders(subg)
				})

			}
		})
	}
}

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
	slices.SortFunc(g.Nodes, func(a, b *graph.Node) bool {
		if a.Layer != b.Layer {
			return a.Layer < b.Layer
		}
		return a.LayerIdx < b.LayerIdx
	})
	for _, n := range g.Nodes {
		fmt.Printf("%s L:%d I:%d\n", n.ID, n.Layer, n.LayerIdx)
	}
}

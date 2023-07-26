package autog

import (
	"testing"

	"github.com/nulab/autog/cyclebreaking"
	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/internal/elk"
	"github.com/nulab/autog/internal/testfiles"
	"github.com/nulab/autog/layering"
	"github.com/nulab/autog/ordering"
	"github.com/nulab/autog/positioning"
)

func TestGansnerNorthOrdering(t *testing.T) {
	testgs := testfiles.ReadTestDir("../internal/testfiles/elk_relabeled")
	for _, g := range testgs[:1] {
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
					ordering.GansnerNorth.Process(subg)
					positioning.VerticalAlign.Process(subg)

				})

			}
		})
	}
}

func setCoords(g *graph.DGraph, elkg *elk.Graph) {
	m := map[string]graph.Size{}
	for _, n := range elkg.Nodes {
		m[n.ID] = graph.Size{H: n.Height, W: n.Width}
	}
	for _, n := range g.Nodes {
		n.H = m[n.ID].H
		n.W = m[n.ID].W
	}
}

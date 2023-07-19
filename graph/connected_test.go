package graph

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vibridi/autog/internal/testutils"
)

func TestConnectedComponents(t *testing.T) {
	tests := []struct {
		filename string
		want     [][]string
	}{
		{
			"algebraic_heateropentank_HeaterOpenTank.json",
			[][]string{{"N1", "N6"}, {"N2", "N10", "N11", "N12"}, {"N3", "N4", "N5", "N7", "N8", "N9"}},
		},
		{
			"ddf_LoopSyntactic2.json",
			[][]string{{"N7", "N8", "N1", "N4", "N3", "N2", "N12", "N5", "N13", "N11", "N6", "N9", "N10"}},
		},
		{
			"gt_diningphilosophers_DiningPhilosophers.json",
			[][]string{{"N1", "N11"}, {"N20"}, {"N21"}, {}, {}},
		},
	}

	for _, c := range tests {
		t.Run(c.filename, func(t *testing.T) {
			g := testutils.ReadTestFile("../testutils/elk/cyclic", c.filename)
			dg := FromAdjacencyList(g.AdjacencyList())
			fmt.Printf("%q\n", nodeIds(dg))

			subgs := dg.ConnectedComponents()
			assertComponents(t, subgs, c.want)
		})
	}
}

func assertComponents(t *testing.T, got []*DGraph, want [][]string) {
	require.Len(t, got, len(want))
	m := map[string][]string{}
	for _, gi := range want {
		populate(m, gi)
	}

	for _, subg := range got {
		assert.ElementsMatch(t, m[subg.Nodes[0].ID], nodeIds(subg))
	}
}

func populate(m map[string][]string, ids []string) {
	for _, id := range ids {
		m[id] = ids
	}
}

func nodeIds(g *DGraph) []string {
	ids := make([]string, len(g.Nodes))
	for i, n := range g.Nodes {
		ids[i] = n.ID
	}
	return ids
}

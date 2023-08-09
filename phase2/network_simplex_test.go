package phase2

import (
	"testing"

	"github.com/nulab/autog/graph"
)

var params = graph.Params{
	NetworkSimplexThoroughness: 28,
	NetworkSimplexBalance:      true,
}

func TestNetworkSimplexLayering(t *testing.T) {
	// testgs := testfiles.ReadTestDir("../internal/testfiles/elk_relabeled")
	// for _, g := range testgs {
	// 	dg := graph.FromElk(g)
	// 	phase1.Greedy.Process(dg, params)
	//
	// 	t.Run(g.Name, func(t *testing.T) {
	// 		for _, subg := range dg.ConnectedComponents() {
	// 			t.Run("component:"+subg.Nodes[0].ID, func(t *testing.T) {
	// 				execNetworkSimplex(subg, graph.Params{})
	// 				for _, e := range subg.Edges {
	// 					assert.True(t, e.From.Layer <= e.To.Layer)
	// 				}
	// 			})
	//
	// 		}
	// 	})
	// }
}

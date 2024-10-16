package autog

import (
	"fmt"
	"testing"

	"github.com/nulab/autog/graph"
	ig "github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
	"github.com/stretchr/testify/assert"
)

var dotAbstract = [][]string{
	{"S24", "27"},
	{"S24", "25"},
	{"S1", "10"},
	{"S1", "2"},
	{"S35", "36"},
	{"S35", "43"},
	{"S30", "31"},
	{"S30", "33"},
	{"9", "42"},
	{"9", "T1"},
	{"25", "T1"},
	{"25", "26"},
	{"27", "T24"},
	{"2", "3"},
	{"2", "16"},
	{"2", "17"},
	{"2", "T1"},
	{"2", "18"},
	{"10", "11"},
	{"10", "14"},
	{"10", "T1"},
	{"10", "13"},
	{"10", "12"},
	{"31", "T1"},
	{"31", "32"},
	{"33", "T30"},
	{"33", "34"},
	{"42", "4"},
	{"26", "4"},
	{"3", "4"},
	{"16", "15"},
	{"17", "19"},
	{"18", "29"},
	{"11", "4"},
	{"14", "15"},
	{"37", "39"},
	{"37", "41"},
	{"37", "38"},
	{"37", "40"},
	{"13", "19"},
	{"12", "29"},
	{"43", "38"},
	{"43", "40"},
	{"36", "19"},
	{"32", "23"},
	{"34", "29"},
	{"39", "15"},
	{"41", "29"},
	{"38", "4"},
	{"40", "19"},
	{"4", "5"},
	{"19", "21"},
	{"19", "20"},
	{"19", "28"},
	{"5", "6"},
	{"5", "T35"},
	{"5", "23"},
	{"21", "22"},
	{"20", "15"},
	{"28", "29"},
	{"6", "7"},
	{"15", "T1"},
	{"22", "23"},
	{"22", "T35"},
	{"29", "T30"},
	{"7", "T8"},
	{"23", "T24"},
	{"23", "T1"},
}

func TestLayoutBugfix(t *testing.T) {
	t.Run("output layout empty nodes", func(t *testing.T) {
		src := graph.EdgeSlice([][]string{
			{"N1", "N2"},
			{"N2", "N3"},
			{"N1", "N3"},
		})
		layout := Layout(
			src,
			WithPositioning(PositioningVAlign),
			WithEdgeRouting(EdgeRoutingNoop),
		)
		assert.Len(t, layout.Nodes, 3)
		assert.Len(t, layout.Edges, 3)
	})

	t.Run("self-loop", func(t *testing.T) {
		t.Run("program halts", func(t *testing.T) {
			src := graph.EdgeSlice([][]string{
				{"a", "b"},
				{"b", "c"},
				{"b", "b"},
				{"a", "d"},
			})
			assert.NotPanics(t, func() { _ = Layout(src) })
		})
		t.Run("successful with single node", func(t *testing.T) {
			src := graph.EdgeSlice([][]string{
				{"a", "a"},
			})
			assert.NotPanics(t, func() { _ = Layout(src) })
		})
	})

	t.Run("greedy cycle breaker fails to break cycles", func(t *testing.T) {
		g := graph.EdgeSlice([][]string{
			{"N1", "N4"},
			{"N1", "N8"},
			{"N2", "N5"},
			{"N2", "N8"},
			{"N3", "N8"},
			{"N6", "N3"},
			{"N8", "N1"},
			{"N8", "N2"},
			{"N8", "N7"},
			{"N8", "N15"},
			{"N8", "N16"},
			{"N9", "N10"},
			{"N10", "N11"},
			{"N12", "N13"},
			{"N13", "N14"},
			{"N15", "N1"},
			{"N15", "N9"},
			{"N15", "N10"},
			{"N16", "N2"},
			{"N16", "N12"},
			{"N16", "N13"},
		})
		assert.NotPanics(t, func() {
			_ = Layout(
				g,
				WithNonDeterministicGreedyCycleBreaker(),
			)
		})
	})

	t.Run("wmedian", func(t *testing.T) {
		t.Run("identical edge segfault in cross counting", func(t *testing.T) {
			src := graph.EdgeSlice([][]string{
				{"gql", "acc"},
				{"gql", "dia"},
				{"gql", "edt"},
				{"gql", "fld"},
				{"gql", "itg"},
				{"gql", "ntf"},
				{"gql", "org"},
				{"gql", "sub"},
				{"gql", "spt"},
				{"gql", "tmp"},
				{"acc", "lgc"},
				{"acc", "sub"},
				{"fld", "acc"},
				{"fld", "dia"},
				{"fld", "org"},
				{"fld", "sub"},
				{"dia", "acc"},
				{"dia", "fld"},
				{"dia", "lgc"},
				{"dia", "org"},
				{"dia", "sub"},
			})
			assert.NotPanics(t, func() { _ = Layout(src) })
		})

		t.Run("wrong initialization of flat edges", func(t *testing.T) {
			src := graph.EdgeSlice([][]string{
				{"gql", "acc"},
				{"gql", "dia"},
				{"gql", "edt"},
				{"gql", "fld"},
				{"gql", "itg"},
				{"gql", "ntf"},
				{"gql", "org"},
				{"gql", "sub"},
				{"gql", "spt"},
				{"gql", "tmp"},
				{"acc", "lgc"},
				{"acc", "sub"},
				{"dia", "acc"},
				{"dia", "fld"},
				{"dia", "lgc"},
				{"dia", "org"},
				{"dia", "sub"},
				{"fld", "acc"},
				{"fld", "dia"},
				{"fld", "org"},
				{"fld", "sub"},
			})
			assert.NotPanics(t, func() { _ = Layout(src) })
		})

		t.Run("wrong handling of fixed positions in wmedian", func(t *testing.T) {
			c := make(chan any, 1)
			assert.NotPanics(t, func() {
				_ = Layout(
					graph.EdgeSlice(dotAbstract),
					WithPositioning(PositioningNoop),
					WithEdgeRouting(EdgeRoutingNoop),
					WithMonitor(imonitor.NewFilteredChan(c, imonitor.MatchAll(3, "gvdot", "crossings"))),
				)
			})

			assert.Equal(t, 46, <-c)
		})
	})

	t.Run("network simplex positioner no panic", func(t *testing.T) {
		src := graph.EdgeSlice(dotAbstract)
		assert.NotPanics(t, func() {
			_ = Layout(
				src,
				WithPositioning(PositioningNetworkSimplex),
				WithEdgeRouting(EdgeRoutingNoop),
				WithNodeFixedSize(100, 100),
			)
		})
	})

	t.Run("b&k no overlaps", func(t *testing.T) {
		g := &ig.DGraph{}
		graph.EdgeSlice([][]string{
			{"a", "b"},
			{"b", "c"},
			{"a", "f"},
			{"f", "g"},
			{"a", "u"},
			{"f", "c"},
			{"c", "k"},
			{"f", "k"},
		}).Populate(g)
		_ = Layout(g,
			WithPositioning(PositioningBrandesKoepf),
			WithNodeFixedSize(130, 60),
		)

		overlaps := 0
		for _, l := range g.Layers {
			for j := 1; j < l.Len(); j++ {
				cur := l.Nodes[j]
				prv := l.Nodes[j-1]

				if prv.X+prv.W > cur.X {
					if overlaps >= 0 {
						// note: this isn't a strict inequality because virtual nodes have size 0x0
						assert.Truef(t, prv.X+prv.W <= cur.X, "%s(X:%.2f) overlaps %s(X+W:%.2f)", cur, cur.X, prv, prv.X+prv.W)
					} else {
						fmt.Printf("warning: overlap between nodes %v and %v within tolerance\n", cur, prv)
					}
					overlaps++
				}
			}
		}
	})

	t.Run("sink coloring program hangs", func(t *testing.T) {
		src := graph.EdgeSlice([][]string{
			{"N1", "N2"},
			{"N3", "N1"},
			{"N2", "N3"},
			{"Nh", "N1"},
			{"Nk", "N1"},
			{"Na", "N2"},
			{"Na", "N3"},
			{"N2", "Nd"},
		})
		assert.NotPanics(t, func() { _ = Layout(src, WithPositioning(PositioningSinkColoring)) })
	})
}

func TestOutputVirtualNodes(t *testing.T) {
	src := graph.EdgeSlice([][]string{
		{"N1", "N2"},
		{"N2", "N3"},
		{"N1", "N3"},
	})
	t.Run("keep virtual nodes", func(t *testing.T) {
		g := &ig.DGraph{}
		src.Populate(g)
		layout := Layout(
			g,
			WithPositioning(PositioningVAlign),
			WithEdgeRouting(EdgeRoutingNoop),
			WithOutputVirtualNodes(true),
		)
		assert.Len(t, layout.Nodes, 4)
		assert.Len(t, layout.Edges, 3)
	})

	t.Run("clip output nodes", func(t *testing.T) {
		g := &ig.DGraph{}
		src.Populate(g)
		layout := Layout(
			g,
			WithPositioning(PositioningVAlign),
			WithEdgeRouting(EdgeRoutingNoop),
			WithOutputVirtualNodes(false),
		)
		assert.Equal(t, 3, len(layout.Nodes))
		assert.Equal(t, 3, cap(layout.Nodes))
		assert.Equal(t, 4, len(g.Nodes))

		assert.Len(t, layout.Edges, 3)
	})
}

func TestNoRegression(t *testing.T) {
	t.Run("ELK", func(t *testing.T) {
		var lib_decg_DECGPi = [][]string{
			{"N2", "N8"},
			{"N2", "N13"},
			{"N2", "N15"},
			{"N2", "N4"},
			{"N3", "N1"},
			{"N4", "N3"},
			{"N5", "N16"},
			{"N5", "N18"},
			{"N6", "N8"},
			{"N6", "N18"},
			{"N7", "N6"},
			{"N8", "N5"},
			{"N8", "N9"},
			{"N9", "N6"},
			{"N9", "N7"},
			{"N10", "N14"},
			{"N10", "N19"},
			{"N11", "N10"},
			{"N12", "N10"},
			{"N12", "N11"},
			{"N13", "N14"},
			{"N14", "N17"},
			{"N14", "N12"},
			{"N15", "N13"},
			{"N16", "N4"},
			{"N17", "N16"},
			{"N17", "N19"},
			{"N18", "N5"},
			{"N19", "N17"},
		}

		var pn_brockackerman_BrockAckerman = [][]string{
			{"N1", "N2"},
			{"N2", "N9"},
			{"N4", "N10"},
			{"N4", "N15"},
			{"N5", "N6"},
			{"N6", "N14"},
			{"N8", "N1"},
			{"N8", "N3"},
			{"N9", "N12"},
			{"N10", "N11"},
			{"N11", "N12"},
			{"N12", "N8"},
			{"N13", "N5"},
			{"N13", "N7"},
			{"N14", "N16"},
			{"N15", "N16"},
			{"N16", "N13"},
		}
		var elkTestGraphs = []struct {
			name string
			adj  [][]string
		}{
			{"lib_decg_DECGPi", lib_decg_DECGPi},
			{"pn_brockackerman_BrockAckerman", pn_brockackerman_BrockAckerman},
		}

		for _, testcase := range elkTestGraphs {
			t.Run(testcase.name, func(t *testing.T) {
				assert.NotPanics(t, func() { Layout(graph.EdgeSlice(testcase.adj)) })
			})
		}
	})

	t.Run("Dot abstract with default options", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Layout(graph.EdgeSlice(dotAbstract))
		})
	})
}

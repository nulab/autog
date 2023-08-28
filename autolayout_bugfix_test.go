package autog

import (
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/monitor"
	"github.com/nulab/autog/phase4"
	"github.com/stretchr/testify/assert"
)

// These tests are meant to ensure a minimal degree of reliability of the default layout pipeline and
// a first barrier against regressions, pending proper unit tests
func TestLayoutCrashers(t *testing.T) {
	t.Run("phase4 SinkColoring", func(t *testing.T) {
		t.Run("#1 and #4", func(t *testing.T) {
			g := graph.FromEdgeSlice([][]string{
				{"N1", "N2"},
				{"N3", "N1"},
				{"N2", "N3"},
				{"Nh", "N1"},
				{"Nk", "N1"},
				{"Na", "N2"},
				{"Na", "N3"},
				{"N2", "Nd"},
			})
			assert.NotPanics(t, func() {
				g = Layout(g, WithPositioning(phase4.SinkColoring))
			})
		})
	})

	t.Run("phase3 GraphvizDot", func(t *testing.T) {
		t.Run("identical edge segfault in cross counting", func(t *testing.T) {
			g := graph.FromEdgeSlice([][]string{
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
			assert.NotPanics(t, func() {
				g = Layout(g)
			})
		})

		t.Run("wrong initialization of flat edges", func(t *testing.T) {
			g := graph.FromEdgeSlice([][]string{
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
			assert.NotPanics(t, func() {
				g = Layout(g)
			})
		})

		t.Run("wrong handling of fixed positions in wmedian", func(t *testing.T) {
			g := graph.FromEdgeSlice([][]string{
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
			})

			c := make(chan monitor.Log, 1)
			assert.NotPanics(t, func() {
				g = Layout(
					g,
					WithPositioning(0),
					WithEdgeRouting(0),
					WithMonitor(monitor.New(c)),
				)
			})

			l := <-c
			assert.Equal(t, 46, l.Value.AsInt())

			want := expectedLayersAbstract()
			for _, n := range g.Nodes {
				if n.IsVirtual {
					continue
				}
				assert.Equalf(t, want[n.ID], n.Layer, "node %s layer %d but should be %d", n.ID, n.Layer, want[n.ID])
			}
		})
	})
}

func expectedLayersAbstract() map[string]int {
	// in dot the nodes 39 and 41 end up inverted
	// this is likely due to a different process order in the vbalance step
	// dot uses qsort which is unstable for equal values
	return map[string]int{
		"S1": 0, "S35": 0,
		"10": 1, "2": 1, "37": 1, "36": 1, "43": 1, "S24": 1,
		"S30": 2, "13": 2, "17": 2, "39": 4, "40": 2, "9": 2, "38": 2, "25": 2,
		"33": 3, "12": 3, "16": 3, "19": 3, "42": 3, "11": 3, "3": 3, "26": 3, "27": 3,
		"34": 4, "18": 4, "41": 2, "28": 4, "31": 4, "14": 4, "20": 4, "21": 4, "4": 4,
		"29": 5, "32": 5, "15": 5, "22": 5, "5": 5,
		"T30": 6, "23": 6, "T35": 6, "6": 6,
		"T1": 7, "T24": 7, "7": 7,
		"T8": 8,
	}
}

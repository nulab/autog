package testfiles

import (
	"testing"

	"github.com/nulab/autog"
	"github.com/nulab/autog/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
	"github.com/nulab/autog/phase4"
	"github.com/stretchr/testify/assert"
)

func TestCrashers(t *testing.T) {
	t.Run("phase4 SinkColoring", func(t *testing.T) {
		t.Run("#1 and #4", func(t *testing.T) {
			g := graph.FromEdgeSlice(issues1and4)
			assert.NotPanics(t, func() { g = autog.Layout(g, autog.WithPositioning(phase4.SinkColoring)) })
		})
	})

	t.Run("phase3 GraphvizDot", func(t *testing.T) {
		t.Run("identical edge segfault in cross counting", func(t *testing.T) {
			g := graph.FromEdgeSlice(cacooArch)
			assert.NotPanics(t, func() { g = autog.Layout(g) })
		})

		t.Run("wrong initialization of flat edges", func(t *testing.T) {
			g := graph.FromEdgeSlice(cacooArch2)
			assert.NotPanics(t, func() { g = autog.Layout(g) })
		})

		t.Run("wrong handling of fixed positions in wmedian", func(t *testing.T) {
			g := graph.FromEdgeSlice(dotAbstract)
			c := make(chan any, 1)
			assert.NotPanics(t, func() {
				g = autog.Layout(
					g,
					autog.WithPositioning(0),
					autog.WithEdgeRouting(0),
					autog.WithMonitor(imonitor.NewFilteredChan(c, imonitor.MatchAll(3, "gvdot", "crossings"))),
				)
			})

			assert.Equal(t, 46, <-c)

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

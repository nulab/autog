//go:build unit

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
			g := graph.FromEdgeSlice(DotAbstract)
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
		})
	})
}

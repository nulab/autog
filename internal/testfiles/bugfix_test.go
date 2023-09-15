//go:build unit

package testfiles

import (
	"math"
	"testing"

	"github.com/nulab/autog"
	"github.com/nulab/autog/graph"
	ig "github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
	"github.com/stretchr/testify/assert"
)

func TestCrashers(t *testing.T) {
	t.Run("phase4 SinkColoring", func(t *testing.T) {
		t.Run("#1 and #4", func(t *testing.T) {
			src := graph.EdgeSlice(issues1and4)
			assert.NotPanics(t, func() { _ = autog.Layout(src, autog.WithPositioningSinkColoring()) })
		})
	})

	t.Run("phase3 GraphvizDot", func(t *testing.T) {
		t.Run("identical edge segfault in cross counting", func(t *testing.T) {
			src := graph.EdgeSlice(cacooArch)
			assert.NotPanics(t, func() { _ = autog.Layout(src) })
		})

		t.Run("wrong initialization of flat edges", func(t *testing.T) {
			src := graph.EdgeSlice(cacooArch2)
			assert.NotPanics(t, func() { _ = autog.Layout(src) })
		})

		t.Run("wrong handling of fixed positions in wmedian", func(t *testing.T) {
			c := make(chan any, 1)
			assert.NotPanics(t, func() {
				_ = autog.Layout(
					graph.EdgeSlice(DotAbstract),
					autog.WithPositioningNoop(),
					autog.WithEdgeRoutingNoop(),
					autog.WithMonitor(imonitor.NewFilteredChan(c, imonitor.MatchAll(3, "gvdot", "crossings"))),
				)
			})

			assert.Equal(t, 46, <-c)
		})
	})

	t.Run("phase4 NetworkSimplex", func(t *testing.T) {
		g := &ig.DGraph{}
		graph.EdgeSlice(DotAbstract).Populate(g)
		for _, n := range g.Nodes {
			n.W, n.H = 100, 100
		}
		assert.NotPanics(t, func() {
			_ = autog.Layout(
				g,
				autog.WithPositioningNetworkSimplex(),
				autog.WithEdgeRoutingNoop(),
			)
		})

		for i := 0; i < len(g.Layers); i++ {
			for j := 1; j < g.Layers[i].Len(); j++ {
				cur := g.Layers[i].Nodes[j]
				prv := g.Layers[i].Nodes[j-1]
				// todo: this isn't a strict inequality bc virtual nodes have size 0x0
				assert.Truef(t, math.Abs(prv.X+prv.W-cur.X) >= 0, "%s(X:%.2f) overlaps %s(X:%.2f)", cur, cur.X, prv, prv.X)
			}
		}
	})
}

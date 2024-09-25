//go:build unit

package testfiles

import (
	"fmt"
	"testing"

	"github.com/nulab/autog"
	"github.com/nulab/autog/graph"
	ig "github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
	"github.com/stretchr/testify/assert"
)

func TestCrashers(t *testing.T) {
	t.Run("phase4 SinkColoring", func(t *testing.T) {
		t.Run("program hangs", func(t *testing.T) {
			src := graph.EdgeSlice(issues1and4)
			assert.NotPanics(t, func() { _ = autog.Layout(src, autog.WithPositioning(autog.PositioningSinkColoring)) })
		})
	})

	t.Run("phase3 WMedian", func(t *testing.T) {
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
					autog.WithPositioning(autog.PositioningNoop),
					autog.WithEdgeRouting(autog.EdgeRoutingNoop),
					autog.WithMonitor(imonitor.NewFilteredChan(c, imonitor.MatchAll(3, "gvdot", "crossings"))),
				)
			})

			assert.Equal(t, 46, <-c)
		})
	})

	t.Run("phase4 NetworkSimplex", func(t *testing.T) {
		g := &ig.DGraph{}
		graph.EdgeSlice(DotAbstract).Populate(g)
		assert.NotPanics(t, func() {
			_ = autog.Layout(
				g,
				autog.WithPositioning(autog.PositioningNetworkSimplex),
				autog.WithEdgeRouting(autog.EdgeRoutingNoop),
				autog.WithNodeFixedSize(100, 100),
			)
		})

		assertNoOverlaps(t, g, 1)
	})

	t.Run("phase4 B&K", func(t *testing.T) {
		t.Run("no overlaps", func(t *testing.T) {
			g := &ig.DGraph{}
			graph.EdgeSlice(bkWrongAlignment).Populate(g)
			_ = autog.Layout(g,
				autog.WithPositioning(autog.PositioningBrandesKoepf),
				autog.WithNodeFixedSize(130, 60),
			)
			
			assertNoOverlaps(t, g, 0)
		})
	})

	t.Run("output layout empty nodes", func(t *testing.T) {
		src := graph.EdgeSlice(simpleVirtualNodes)
		layout := autog.Layout(
			src,
			autog.WithPositioning(autog.PositioningVAlign),
			autog.WithEdgeRouting(autog.EdgeRoutingNoop),
		)
		assert.Len(t, layout.Nodes, 3)
		assert.Len(t, layout.Edges, 3)
	})

	t.Run("self-loop", func(t *testing.T) {
		t.Run("program halts", func(t *testing.T) {
			src := graph.EdgeSlice(withSelfLoop)
			assert.NotPanics(t, func() { _ = autog.Layout(src) })
		})
		t.Run("successful with single node", func(t *testing.T) {
			src := graph.EdgeSlice(singleSelfLoop)
			assert.NotPanics(t, func() { _ = autog.Layout(src) })
		})
	})
}

func assertNoOverlaps(t *testing.T, g *ig.DGraph, tolerance int) {
	overlaps := 0
	for i := 0; i < len(g.Layers); i++ {
		for j := 1; j < g.Layers[i].Len(); j++ {
			cur := g.Layers[i].Nodes[j]
			prv := g.Layers[i].Nodes[j-1]

			if prv.X+prv.W > cur.X {
				if overlaps >= tolerance {
					// note: this isn't a strict inequality because virtual nodes have size 0x0
					assert.Truef(t, prv.X+prv.W <= cur.X, "%s(X:%.2f) overlaps %s(X+W:%.2f)", cur, cur.X, prv, prv.X+prv.W)
				} else {
					fmt.Printf("warning: overlap between nodes %v and %v within tolerance\n", cur, prv)
				}
				overlaps++
			}
		}
	}
}

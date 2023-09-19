package autog

import (
	"testing"

	"github.com/nulab/autog/internal/graph"
	"github.com/nulab/autog/internal/phase1"
	"github.com/nulab/autog/internal/phase2"
	"github.com/nulab/autog/internal/phase3"
	"github.com/nulab/autog/internal/phase4"
	"github.com/nulab/autog/internal/phase5"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	// set some random options that are different from the defaults
	opts := testOptions(
		WithCycleBreaking(CycleBreakingDepthFirst),
		WithOrdering(OrderingNoop),
		WithPositioning(PositioningVAlign),
		WithEdgeRouting(EdgeRoutingStraight),
		WithNetworkSimplexThoroughness(30),
		WithLayerSpacing(75.5),
		WithNodeSpacing(10.0),
		WithBrandesKoepfLayout(2),
	)

	assert.Equal(t, phase1.DepthFirst, opts.p1)
	assert.Equal(t, phase3.NoOrdering, opts.p3)
	assert.Equal(t, phase4.VerticalAlign, opts.p4)
	assert.Equal(t, phase5.Straight, opts.p5)
	assert.Equal(t, uint(30), opts.params.NetworkSimplexThoroughness)
	assert.Equal(t, graph.OptionNsBalanceV, opts.params.NetworkSimplexBalance)
	assert.Equal(t, 75.5, opts.params.LayerSpacing)
	assert.Equal(t, 10.0, opts.params.NodeSpacing)
	assert.Equal(t, 2, opts.params.BrandesKoepfLayout)
	assert.Nil(t, opts.monitor)

	assert.Equal(t, CycleBreakingGreedy, phase1.Greedy)
	assert.Equal(t, CycleBreakingDepthFirst, phase1.DepthFirst)
	assert.Equal(t, LayeringLongestPath, phase2.LongestPath)
	assert.Equal(t, LayeringNetworkSimplex, phase2.NetworkSimplex)
	assert.Equal(t, OrderingNoop, phase3.NoOrdering)
	assert.Equal(t, OrderingWMedian, phase3.WMedian)
	assert.Equal(t, PositioningNoop, phase4.NoPositioning)
	assert.Equal(t, PositioningVAlign, phase4.VerticalAlign)
	assert.Equal(t, PositioningNetworkSimplex, phase4.NetworkSimplex)
	assert.Equal(t, PositioningBrandesKoepf, phase4.BrandesKoepf)
	assert.Equal(t, PositioningSinkColoring, phase4.SinkColoring)
	assert.Equal(t, PositioningPackRight, phase4.PackRight)
	assert.Equal(t, EdgeRoutingNoop, phase5.NoRouting)
	assert.Equal(t, EdgeRoutingStraight, phase5.Straight)
	assert.Equal(t, EdgeRoutingPolyline, phase5.Polyline)
	assert.Equal(t, EdgeRoutingOrtho, phase5.Ortho)
}

func testOptions(opts ...Option) options {
	layoutOpts := defaultOptions
	for _, opt := range opts {
		opt(&layoutOpts)
	}
	return layoutOpts
}

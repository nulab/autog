package autog

import (
	"testing"

	"github.com/nulab/autog/internal/graph"
	"github.com/nulab/autog/internal/phase1"
	"github.com/nulab/autog/internal/phase3"
	"github.com/nulab/autog/internal/phase4"
	"github.com/nulab/autog/internal/phase5"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	// set some random options that are different from the defaults
	opts := testOptions(
		WithCycleBreakingDFS(),
		WithOrderingNoop(),
		WithPositioningVAlign(),
		WithEdgeRoutingStraight(),
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
}

func testOptions(opts ...Option) options {
	layoutOpts := defaultOptions
	for _, opt := range opts {
		opt(&layoutOpts)
	}
	return layoutOpts
}

package autog

import (
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/phase1"
	"github.com/nulab/autog/phase2"
	"github.com/nulab/autog/phase3"
	"github.com/nulab/autog/phase4"
	"github.com/nulab/autog/phase5"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	// set some random options that are different from the defaults
	opts := testOptions(
		WithCycleBreaking(phase1.DepthFirst),
		WithLayering(phase2.Alg(12)),
		WithOrdering(phase3.NoOrdering),
		WithPositioning(phase4.VerticalAlign),
		WithEdgeRouting(phase5.Straight),
		WithNetworkSimplexThoroughness(30),
		WithNetworkSimplexBalance(graph.OptionNsBalanceH),
		WithLayerSpacing(75.5),
		WithNodeSpacing(10.0),
		WithBrandesKoepfLayout(2),
	)

	assert.Equal(t, phase1.DepthFirst, opts.p1)
	assert.Equal(t, phase2.Alg(12), opts.p2)
	assert.Equal(t, phase3.NoOrdering, opts.p3)
	assert.Equal(t, phase4.VerticalAlign, opts.p4)
	assert.Equal(t, phase5.Straight, opts.p5)
	assert.Equal(t, uint(30), opts.params.NetworkSimplexThoroughness)
	assert.Equal(t, graph.OptionNsBalanceH, opts.params.NetworkSimplexBalance)
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

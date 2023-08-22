package autog

import (
	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/monitor"
	cbreaking "github.com/nulab/autog/phase1"
	layering "github.com/nulab/autog/phase2"
	ordering "github.com/nulab/autog/phase3"
	positioning "github.com/nulab/autog/phase4"
	routing "github.com/nulab/autog/phase5"
)

type options struct {
	p1     cbreaking.Alg
	p2     layering.Alg
	p3     ordering.Alg
	p4     positioning.Alg
	p5     routing.Alg
	params graph.Params
}

var defaultOptions = options{
	p1: cbreaking.Greedy,
	p2: layering.NetworkSimplex,
	p3: ordering.GraphvizDot,
	p4: positioning.SinkColoring,
	p5: routing.PieceWise,
	params: graph.Params{
		NetworkSimplexThoroughness:               28,
		NetworkSimplexMaxIterFactor:              0,
		NetworkSimplexBalance:                    true,
		GraphvizDotMaxIter:                       24,
		NetworkSimplexAuxiliaryGraphWeightFactor: 4,
		LayerSpacing:                             150.0,
		NodeSpacing:                              60.0,
		BrandesKoepfLayout:                       -1,
		Monitor:                                  monitor.NewNoop(),
	},
}

type Option func(*options)

func WithCycleBreaking(alg cbreaking.Alg) Option {
	return func(o *options) {
		o.p1 = alg
	}
}

func WithLayering(alg layering.Alg) Option {
	return func(o *options) {
		o.p2 = alg
	}
}

func WithOrdering(alg ordering.Alg) Option {
	return func(o *options) {
		o.p3 = alg
	}
}

func WithPositioning(alg positioning.Alg) Option {
	return func(o *options) {
		o.p4 = alg
	}
}

func WithEdgeRouting(alg routing.Alg) Option {
	return func(o *options) {
		o.p5 = alg
	}
}

func WithNetworkSimplexThoroughness(thoroughness uint) Option {
	return func(o *options) {
		o.params.NetworkSimplexThoroughness = thoroughness
	}
}

func WithNetworkSimplexBalance(balance bool) Option {
	return func(o *options) {
		o.params.NetworkSimplexBalance = balance
	}
}

func WithLayerSpacing(spacing float64) Option {
	return func(o *options) {
		o.params.LayerSpacing = spacing
	}
}

func WithNodeSpacing(spacing float64) Option {
	return func(o *options) {
		o.params.NodeSpacing = spacing
	}
}

func WithBrandesKoepfLayout(i int) Option {
	return func(o *options) {
		o.params.BrandesKoepfLayout = i
	}
}

func WithMonitor(monitor monitor.Monitor) Option {
	return func(o *options) {
		o.params.Monitor = monitor
	}
}

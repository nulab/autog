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
	p5: routing.NoRouting,
	params: graph.Params{
		NetworkSimplexThoroughness:               28,
		NetworkSimplexMaxIterFactor:              0,
		NetworkSimplexBalance:                    true,
		GraphvizDotMaxIter:                       24,
		NetworkSimplexAuxiliaryGraphWeightFactor: 4,
		LayerSpacing:                             150.0,
		NodeMargin:                               20.0,
		NodeSpacing:                              40.0,
		Monitor:                                  monitor.NewNoop(),
	},
}

type option func(*options)

func WithCycleBreaking(alg cbreaking.Alg) option {
	return func(o *options) {
		o.p1 = alg
	}
}

func WithLayering(alg layering.Alg) option {
	return func(o *options) {
		o.p2 = alg
	}
}

func WithOrdering(alg ordering.Alg) option {
	return func(o *options) {
		o.p3 = alg
	}
}

func WithPositioning(alg positioning.Alg) option {
	return func(o *options) {
		o.p4 = alg
	}
}

func WithEdgeRouting(alg routing.Alg) option {
	return func(o *options) {
		o.p5 = alg
	}
}

func WithMonitor(monitor monitor.Monitor) option {
	return func(o *options) {
		o.params.Monitor = monitor
	}
}

// todo: param options

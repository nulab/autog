package autog

import (
	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
	cbreaking "github.com/nulab/autog/internal/phase1"
	layering "github.com/nulab/autog/internal/phase2"
	ordering "github.com/nulab/autog/internal/phase3"
	positioning "github.com/nulab/autog/internal/phase4"
	routing "github.com/nulab/autog/internal/phase5"
)

type options struct {
	p1      cbreaking.Alg
	p2      layering.Alg
	p3      ordering.Alg
	p4      positioning.Alg
	p5      routing.Alg
	params  graph.Params
	monitor imonitor.Monitor
	output  output
}

type output struct {
	keepVirtualNodes bool
}

var defaultOptions = options{
	p1: cbreaking.Greedy,
	p2: layering.NetworkSimplex,
	p3: ordering.WMedian,
	p4: positioning.SinkColoring,
	p5: routing.Polyline,
	params: graph.Params{
		NetworkSimplexThoroughness:               28,
		NetworkSimplexMaxIterFactor:              0,
		NetworkSimplexBalance:                    graph.OptionNsBalanceV,
		WMedianMaxIter:                           24,
		NetworkSimplexAuxiliaryGraphWeightFactor: 4,
		LayerSpacing:                             150.0,
		NodeSpacing:                              60.0,
		BrandesKoepfLayout:                       -1,
	},
	monitor: nil,
	output:  defaultOutputOptions,
}

var defaultOutputOptions = output{
	keepVirtualNodes: false,
}

type Option func(*options)

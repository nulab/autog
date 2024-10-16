package autog

import (
	"github.com/nulab/autog/graph"
	ig "github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

func WithNonDeterministicGreedyCycleBreaker() Option {
	return func(o *options) {
		o.params.GreedyCycleBreakerRandomNodeChoice = true
	}
}

func WithNetworkSimplexThoroughness(thoroughness uint) Option {
	return func(o *options) {
		o.params.NetworkSimplexThoroughness = thoroughness
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

// WithNodeSize sets a size to each node found in the supplied map. The map keys are the node ids.
// Individual node sizes override the size set by WithNodeFixedSize.
func WithNodeSize(sizes map[string]graph.Size) Option {
	return func(o *options) {
		o.params.NodeSizeFunc = func(n *ig.Node) {
			n.Size = sizes[n.ID]
		}
	}
}

// WithNodeFixedSize sets the same size to all nodes in the source graph.
func WithNodeFixedSize(w, h float64) Option {
	return func(o *options) {
		o.params.NodeFixedSizeFunc = func(n *ig.Node) {
			n.W = w
			n.H = h
		}
	}
}

func WithVirtualNodeFixedSize(n float64) Option {
	return func(o *options) {
		o.params.VirtualNodeFixedSize = n
	}
}

func WithBrandesKoepfLayout(i int) Option {
	return func(o *options) {
		o.params.BrandesKoepfLayout = i
	}
}

func WithMonitor(monitor imonitor.Monitor) Option {
	return func(o *options) {
		o.monitor = monitor
	}
}

func WithOutputVirtualNodes(keep bool) Option {
	return func(o *options) {
		o.output.includeVirtual = keep
	}
}

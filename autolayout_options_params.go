package autog

import (
	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

func WithNetworkSimplexThoroughness(thoroughness uint) Option {
	return func(o *options) {
		o.params.NetworkSimplexThoroughness = thoroughness
	}
}

func WithNetworkSimplexBalance(balance graph.OptionNsBalance) Option {
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

func WithMonitor(monitor imonitor.Monitor) Option {
	return func(o *options) {
		o.monitor = monitor
	}
}

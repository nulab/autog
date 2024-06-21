package autog

import (
	ig "github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

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

func WithNodeFixedSize(w, h float64) Option {
	return func(o *options) {
		o.params.NodeFixedSizeFunc = func(n *ig.Node) {
			n.W = w
			n.H = h
		}
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

func WithKeepVirtualNodes(keep bool) Option {
	return func(o *options) {
		o.output.keepVirtualNodes = keep
	}
}

package autog

import (
	"github.com/nulab/autog/cyclebreaking"
	"github.com/nulab/autog/layering"
	"github.com/nulab/autog/ordering"
	"github.com/nulab/autog/positioning"
)

type options struct {
	p1 cyclebreaking.Alg
	p2 layering.Alg
	p3 ordering.Alg
	p4 positioning.Alg
}

var defaultOptions = options{
	p1: cyclebreaking.Greedy,
	p2: layering.NetworkSimplex,
	p3: ordering.GansnerNorth,
	p4: positioning.VerticalAlign,
}

type option func(*options)

func WithCycleBreaking(alg cyclebreaking.Alg) option {
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

package autog

import (
	"github.com/vibridi/autog/cyclebreaking"
	"github.com/vibridi/autog/layering"
	"github.com/vibridi/autog/ordering"
)

type options struct {
	p1 cyclebreaking.Alg
	p2 layering.Alg
	p3 ordering.Alg
}

var defaultOptions = options{
	p1: cyclebreaking.Greedy,
	p2: layering.NetworkSimplex,
	p3: ordering.GansnerNorth,
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

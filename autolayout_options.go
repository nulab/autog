package autog

import (
	"github.com/vibridi/autog/cyclebreaking"
	"github.com/vibridi/autog/layering"
)

type options struct {
	p1 cyclebreaking.Alg
	p2 layering.Alg
}

var defaultOptions = options{
	p1: cyclebreaking.Greedy,
	p2: layering.NetworkSimplex,
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

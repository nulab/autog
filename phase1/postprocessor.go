package phase1

import (
	"github.com/nulab/autog/graph"
)

type postProcessor uint8

const (
	RestoreEdges postProcessor = iota
)

func (postProcessor) Process(g *graph.DGraph, _ graph.Params) {
	for _, e := range g.Edges {
		if e.IsReversed {
			// reverse back
			e.Reverse()
		}
	}
}

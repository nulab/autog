package phase2

import (
	"github.com/nulab/autog/internal/graph"
	"github.com/nulab/autog/internal/ns"
)

func execNetworkSimplex(g *graph.DGraph, params graph.Params) {
	new(ns.Processor).Exec(g, params)
}

package processor

import (
	"github.com/nulab/autog/internal/graph"
)

// P represents a pipeline processor
type P interface {
	Phase() int
	String() string
	Process(*graph.DGraph, graph.Params)
}

// F represents a standalone function that can execute a processing task.
type F func(*graph.DGraph)

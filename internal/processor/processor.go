package processor

import (
	"github.com/nulab/autog/internal/graph"
)

type P interface {
	Phase() int
	String() string
	Process(*graph.DGraph, graph.Params)
}

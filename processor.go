package autog

import (
	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/monitor"
)

type processor interface {
	Process(*graph.DGraph, *monitor.Monitor)
}

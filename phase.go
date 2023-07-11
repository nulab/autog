package autog

import "github.com/vibridi/autog/internal/graph"

type phase interface {
	Process(g *graph.DGraph)
	Cleanup()
}

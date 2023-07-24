package cyclebreaking

import "github.com/vibridi/autog/graph"

type postProcessor uint8

const (
	UndoRevertEdges postProcessor = iota
)

func (postProcessor) Process(g *graph.DGraph) {
	if !g.HasCycles() {
		return
	}
	for _, e := range g.Edges {
		if e.IsReversed {
			e.IsReversed = false
			e.Reverse()
		}
	}
}

package preprocessor

import (
	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
	"github.com/nulab/autog/internal/processor"
)

func IgnoreSelfLoops(g *graph.DGraph) processor.F {
	del := graph.EdgeSet{}
	for _, e := range g.Edges {
		if e.From == e.To {
			imonitor.Log(imonitor.KeySelfLoop, "removed: "+e.From.ID)
			del[e] = true
		}
	}
	for e := range del {
		// using the appropriate From-Out and To-In fields for clarity,
		// but it's always the same node with the same incoming and outgoing edge
		e.From.Out.Remove(e)
		e.To.In.Remove(e)
		g.Edges.Remove(e)
	}
	return func(g *graph.DGraph) {
		for e := range del {
			imonitor.Log(imonitor.KeySelfLoop, "added: "+e.From.ID)
			e.From.Out.Add(e)
			e.To.In.Add(e)
			g.Edges.Add(e)
		}
	}
}

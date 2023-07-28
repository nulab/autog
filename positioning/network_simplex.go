package positioning

import (
	"strconv"

	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/layering"
)

const (
	edgeWeightFactor = 4
)

type networkSimplexProcessor struct {
	nodes map[string]*graph.Node
}

func execNetworkSimplex(g *graph.DGraph) {
	p := &networkSimplexProcessor{
		nodes: map[string]*graph.Node{},
	}

	aux := p.auxiliaryGraph(g)
	layering.NetworkSimplex.Process(aux)

	for _, l := range g.Layers {
		for i, n := range l.Nodes {
			l.W += defaultNodeMargin*2 + n.W
			l.H = max(l.H, n.H)
			n.X = float64(p.nodes[n.ID].Layer + ((defaultNodeMargin + defaultNodeSpacing) * i))
		}
	}
}

func (p *networkSimplexProcessor) auxiliaryGraph(g *graph.DGraph) *graph.DGraph {
	g1 := &graph.DGraph{}

	for _, n := range g.Nodes {
		m := &graph.Node{ID: n.ID}
		p.nodes[m.ID] = m
		g1.Nodes = append(g1.Nodes, m)
	}
	for i, e := range g.Edges {
		ne := &graph.Node{ID: "NE" + strconv.Itoa(i)}
		p.nodes[ne.ID] = ne
		g1.Nodes = append(g1.Nodes, ne)

		weight := e.Weight * omega(e)

		u, v := p.nodes[e.From.ID], p.nodes[e.To.ID]

		eu := graph.NewEdge(ne, u, weight)
		eu.Delta = 0
		ne.Out = append(ne.Out, eu)
		u.In = append(u.In, eu)

		ev := graph.NewEdge(ne, v, weight)
		ev.Delta = 0
		ne.Out = append(ne.Out, ev)
		v.In = append(v.In, ev)

		g1.Edges = append(g1.Edges, eu, ev)
	}
	for _, l := range g.Layers {
		for i := 0; i < len(l.Nodes)-1; i++ {
			v := p.nodes[l.Nodes[i].ID]
			w := p.nodes[l.Nodes[i+1].ID]
			f := graph.NewEdge(v, w, 0)
			f.Delta = int(rho(v, w)) // probably not correct, conversion to int amounts to a floor()
			g1.Edges = append(g1.Edges, f)

			v.Out = append(v.Out, f)
			w.In = append(w.In, f)
		}
	}
	return g1
}

func rho(a, b *graph.Node) float64 {
	return (a.W+b.W)/2 + 100 // 100 should be default node spacing?
}

// todo this could be merged with Edge.Type
func omega(e *graph.Edge) int {
	switch e.Type() {
	case 0:
		return 1 * edgeWeightFactor
	case 1:
		return 2 * edgeWeightFactor
	case 2:
		return 8 * edgeWeightFactor
	default:
		panic("unexpected edge type")
	}
}

package phase4

import (
	"math"
	"strconv"

	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/phase2"
)

type networkSimplexProcessor struct {
	nodes            map[string]*graph.Node
	edgeWeightFactor int
	nodeSpacing      float64
}

// Positioning algorithm used in Graphviz Dot and described in:
//   - "Emden R. Gansner, Eleftherios Koutsofios, Stephen C. North, Kiem-Phong Vo, A technique for
//     drawing directed graphs. Software Engineering 19(3), pp. 214-230, 1993."
//     https://www.researchgate.net/publication/3187542_A_Technique_for_Drawing_Directed_Graphs
//
// It constructs an auxiliary graph and runs network simplex on it. The resulting layers are the X coordinates of the main graph.
func execNetworkSimplex(g *graph.DGraph, params graph.Params) {
	p := &networkSimplexProcessor{
		nodes:            map[string]*graph.Node{},
		edgeWeightFactor: params.NetworkSimplexAuxiliaryGraphWeightFactor,
		nodeSpacing:      params.NodeSpacing,
	}

	// todo: if there are flat edges, dot adds auxiliary edges
	aux := p.auxiliaryGraph(g)

	params.NetworkSimplexMaxIterFactor = len(g.Nodes)
	params.NetworkSimplexBalance = 2
	phase2.NetworkSimplex.Process(aux, params)

	for _, l := range g.Layers {
		for _, n := range l.Nodes {
			l.H = max(l.H, n.H)
			n.X = float64(p.nodes[n.ID].Layer)
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

		weight := e.Weight * omega(e) * p.edgeWeightFactor

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
			f.Delta = int(math.Round(p.distCenterPoints(v, w)))
			g1.Edges = append(g1.Edges, f)

			v.Out = append(v.Out, f)
			w.In = append(w.In, f)
		}
	}
	return g1
}

func (p *networkSimplexProcessor) distCenterPoints(a, b *graph.Node) float64 {
	return (a.W / 2) + (b.W / 2) + p.nodeSpacing
}

// todo this could be merged with Edge.Type
func omega(e *graph.Edge) int {
	switch e.Type() {
	case 0:
		return 1
	case 1:
		return 2
	case 2:
		return 8
	default:
		panic("unexpected edge type")
	}
}

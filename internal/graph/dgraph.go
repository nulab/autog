package graph

import (
	"strings"
)

type DGraph struct {
	Nodes  []*Node
	Edges  EdgeList
	Layers Layers
}

// Generate implements graph.Source by returning itself as a compatibility step
func (g *DGraph) Generate() *DGraph {
	return g
}

// todo: sources and sinks don't yet account for isolated nodes with a self-loop

// Sources returns a list of nodes with no incoming edges
func (g *DGraph) Sources() []*Node {
	var sources []*Node
	for _, n := range g.Nodes {
		if len(n.In) == 0 {
			sources = append(sources, n)
		}
	}
	return sources
}

// Sinks returns a list of nodes with no outgoing edges
func (g *DGraph) Sinks() []*Node {
	var sinks []*Node
	for _, n := range g.Nodes {
		if len(n.Out) == 0 {
			sinks = append(sinks, n)
		}
	}
	return sinks
}

func (g *DGraph) String() string {
	bld := strings.Builder{}
	for _, n := range g.Nodes {
		bld.WriteString(n.ID)
		bld.WriteRune('\n')
		bld.WriteString("-IN:")
		if len(n.In) == 0 {
			bld.WriteRune('\t')
			bld.WriteString("none")
			bld.WriteRune('\n')
		}
		for _, e := range n.In {
			bld.WriteRune('\t')
			bld.WriteString(e.From.ID)
			bld.WriteString(" -> ")
			bld.WriteString(n.ID)
			bld.WriteRune('\n')
		}
		bld.WriteString("-OUT:")
		if len(n.Out) == 0 {
			bld.WriteRune('\t')
			bld.WriteString("none")
			bld.WriteRune('\n')
		}
		for _, e := range n.Out {
			bld.WriteRune('\t')
			bld.WriteString(n.ID)
			bld.WriteString(" -> ")
			bld.WriteString(e.To.ID)
			bld.WriteRune('\n')
		}
	}
	return bld.String()
}

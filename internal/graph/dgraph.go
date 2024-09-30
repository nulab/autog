package graph

import (
	"iter"
	"strings"
)

type DGraph struct {
	Nodes  []*Node
	Edges  EdgeList
	Layers Layers
}

func (g *DGraph) Populate(*DGraph) {
	// DGraph implements graph.Source to facilitate unit testing.
	// It is not meant to do anything because *DGraph can be obtained from graph.Source
	// via type assertion.
	// todo: instead of having a dummy implementation, a wrapper could be used instead in unit tests
}

func (g *DGraph) GetNodes() []*Node {
	return nil
}

func (g *DGraph) GetEdges() []*Edge {
	return nil
}

// Sources returns a sequence of nodes with no incoming edges
func (g *DGraph) Sources() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for _, n := range g.Nodes {
			if n.Indeg() == 0 {
				if !yield(n) {
					return
				}
			}
		}
	}
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

package graph

import (
	"iter"
)

type DGraph struct {
	Nodes  []*Node
	Edges  EdgeList
	Layers []*Layer
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
func (g *DGraph) Sinks() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for _, n := range g.Nodes {
			if n.Outdeg() == 0 {
				if !yield(n) {
					return
				}
			}
		}
	}
}

func (g *DGraph) VirtualNodes() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for _, n := range g.Nodes {
			if n.IsVirtual {
				if !yield(n) {
					return
				}
			}
		}
	}
}

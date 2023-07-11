package graph

// todo: rename to DGraph to make it clear it's directed
type Graph struct {
	Edges    []*Edge
	Nodes    []*Node
	IsCyclic bool
}

// Sources returns a list of nodes with no incoming edges
func (g *Graph) Sources() []*Node {
	var sources []*Node
	for _, n := range g.Nodes {
		if len(n.In) == 0 {
			sources = append(sources, n)
		}
	}
	return sources
}

// Sinks returns a list of nodes with no outgoing edges
func (g *Graph) Sinks() []*Node {
	var sinks []*Node
	for _, n := range g.Nodes {
		if len(n.Out) == 0 {
			sinks = append(sinks, n)
		}
	}
	return sinks
}

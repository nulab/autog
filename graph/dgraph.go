package graph

import (
	"maps"
	"strings"
)

type DGraph struct {
	Nodes    []*Node
	Edges    []*Edge
	Layers   map[int][]*Node
	LayersX  map[int]*Layer
	isCyclic *bool
}

// todo: this can probably become generic, to allow arbitrary ID types
func FromAdjacencyList(list map[string][]string) *DGraph {
	nodeMap := map[string]*Node{}
	edgeList := []*Edge{}
	for sourceId, targetIds := range list {
		n := nodeMap[sourceId]
		if n == nil {
			n = &Node{ID: sourceId}
		}
		for _, targetId := range targetIds {
			m := nodeMap[targetId]
			if m == nil {
				m = &Node{ID: targetId}
			}
			e := &Edge{From: n, To: m}
			edgeList = append(edgeList, e)

			m.In = append(m.In, e)
			n.Out = append(n.Out, e)

			nodeMap[targetId] = m
		}
		nodeMap[sourceId] = n
	}
	return &DGraph{
		Nodes: maps.Values(nodeMap), // note: this is non-deterministic
		Edges: edgeList,
	}
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
		for _, e := range n.Out {
			bld.WriteString(n.ID)
			bld.WriteString(" -> ")
			bld.WriteString(e.To.ID)
			bld.WriteRune('\n')
		}
	}
	return bld.String()
}

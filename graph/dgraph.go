package graph

import (
	"strings"

	"github.com/nulab/autog/internal/elk"
)

type DGraph struct {
	Nodes       []*Node
	Edges       EdgeList
	HiddenEdges EdgeList
	Layers      Layers
	isCyclic    *bool
}

func FromElk(g *elk.Graph) *DGraph {
	nodeMap := map[string]*Node{}
	portNodeMap := make(map[string]string) // port-node map (each port belongs to one node)

	nodeList := []*Node{}
	edgeList := []*Edge{}

	for _, n := range g.Nodes {
		for _, p := range n.Ports {
			portNodeMap[p.ID] = n.ID
		}
	}
	for _, edge := range g.Edges {
		if len(edge.Sources) > 1 || len(edge.Targets) > 1 {
			panic("hyperedges are not supported")
		}
		sourceId := portNodeMap[edge.Sources[0]]
		targetId := portNodeMap[edge.Targets[0]]

		sourceNode := nodeMap[sourceId]
		if sourceNode == nil {
			sourceNode = &Node{ID: sourceId}
			nodeList = append(nodeList, sourceNode)
			nodeMap[sourceId] = sourceNode
		}
		targetNode := nodeMap[targetId]
		if targetNode == nil {
			targetNode = &Node{ID: targetId}
			nodeList = append(nodeList, targetNode)
			nodeMap[targetId] = targetNode
		}

		e := NewEdge(sourceNode, targetNode, 1)
		edgeList = append(edgeList, e)

		targetNode.In = append(targetNode.In, e)
		sourceNode.Out = append(sourceNode.Out, e)
	}

	return &DGraph{
		Nodes:       nodeList,
		Edges:       edgeList,
		HiddenEdges: EdgeList{},
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

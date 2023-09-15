package graph

import "github.com/nulab/autog/internal/graph"

// EdgeSlice is a graph Source.
type EdgeSlice [][]string

var _ Source = EdgeSlice{}

func (edges EdgeSlice) Generate() *graph.DGraph {
	nodeMap := map[string]*graph.Node{}

	nodeList := []*graph.Node{}
	edgeList := []*graph.Edge{}

	for _, e := range edges {
		if len(e) != 2 {
			panic("graph init: edge must have one source and one target node")
		}
		sourceId := e[0]
		targetId := e[1]

		sourceNode := nodeMap[sourceId]
		if sourceNode == nil {
			sourceNode = &graph.Node{ID: sourceId}
			nodeList = append(nodeList, sourceNode)
			nodeMap[sourceId] = sourceNode
		}
		targetNode := nodeMap[targetId]
		if targetNode == nil {
			targetNode = &graph.Node{ID: targetId}
			nodeList = append(nodeList, targetNode)
			nodeMap[targetId] = targetNode
		}

		e := graph.NewEdge(sourceNode, targetNode, 1)
		edgeList = append(edgeList, e)

		targetNode.In = append(targetNode.In, e)
		sourceNode.Out = append(sourceNode.Out, e)
	}

	return &graph.DGraph{
		Nodes: nodeList,
		Edges: edgeList,
	}
}

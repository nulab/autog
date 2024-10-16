package graph

// EdgeSlice is a graph Source.
type EdgeSlice [][]string

var _ Source = EdgeSlice{}

func (edges EdgeSlice) Populate(g *DGraph) {
	nodeMap := map[string]*Node{}

	nodeList := []*Node{}
	edgeList := []*Edge{}

	for _, e := range edges {
		if len(e) != 2 {
			panic("graph source: edge must have one source and one target node")
		}
		sourceId := e[0]
		targetId := e[1]

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

		e := NewEdge(sourceNode, targetNode, 1) // default to weight 1
		edgeList = append(edgeList, e)

		targetNode.In = append(targetNode.In, e)
		sourceNode.Out = append(sourceNode.Out, e)
	}

	g.Nodes = nodeList
	g.Edges = edgeList
}

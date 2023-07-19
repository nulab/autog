package elk

type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Graph struct {
	ID    string  `json:"id"`
	Name  string  `json:"-"`
	Nodes []*Node `json:"children"`
	Edges []*Edge `json:"edges"`
}

type Node struct {
	ID         string   `json:"id"`
	Ports      []*Port  `json:"ports"`
	Labels     []*Label `json:"labels"`
	Properties *struct {
		PortConstraints     string `json:"portConstraints"`
		NodeLabelsPlacement string `json:"nodeLabels.placement"`
	} `json:"properties"`
	Size
}

type Edge struct {
	ID      string   `json:"id"`
	Sources []string `json:"sources"`
	Targets []string `json:"targets"`
}

type Port struct {
	ID         string `json:"id"`
	Properties *struct {
		Side  string `json:"port.side"`
		Index string `json:"port.index"`
	} `json:"properties"`
}

type Label struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Size
}

func (g *Graph) AdjacencyList() map[string][]string {
	list := make(map[string][]string)
	portNodeMap := make(map[string]string) // port-node map (each port belongs to one node)

	for _, n := range g.Nodes {
		list[n.ID] = []string{}
		for _, p := range n.Ports {
			portNodeMap[p.ID] = n.ID
		}
	}
	for _, edge := range g.Edges {
		if len(edge.Sources) > 1 {
			panic("hyperedges are not supported")
		}
		sourceNode := portNodeMap[edge.Sources[0]]

		targetNodes := list[sourceNode]
		for _, targetPort := range edge.Targets {
			targetNodes = append(targetNodes, portNodeMap[targetPort])
		}
		list[sourceNode] = targetNodes
	}
	return list
}
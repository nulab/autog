package graph

type cycleDetector struct {
	visited  NodeSet
	finished NodeSet
	isCyclic bool
}

func (g *DGraph) HasCycles() bool {
	if g.isCyclic != nil {
		return *g.isCyclic
	}
	d := cycleDetector{
		visited:  make(NodeSet),
		finished: make(NodeSet),
	}
	g.isCyclic = new(bool)

	for _, n := range g.Nodes {
		if !d.visited[n] && !d.finished[n] {
			d.visit(n)
			if d.isCyclic {
				break
			}
		}
	}
	*g.isCyclic = d.isCyclic
	return *g.isCyclic
}

func (d *cycleDetector) visit(n *Node) {
	d.visited[n] = true
	for _, m := range n.AdjacentNodes() {
		if n == m {
			continue // ignore self loops
		}
		if d.visited[m] {
			d.isCyclic = true
			return
		}
		if !d.finished[m] {
			d.visit(m)
		}
	}
	d.visited[n] = false
	d.finished[n] = true
}

package cyclebreaking

import "github.com/vibridi/autog/internal/graph"

type depthFirst struct {
	visited    graph.NodeSet
	active     graph.NodeSet
	reversable []*graph.Edge
}

var DepthFirst = &depthFirst{}

// Process implements a depth-first cycle breaker.
// This is a port of ELK Java code, slightly modified to adapt to Go data structures and coding practices.
// https://github.com/eclipse/elk/blob/master/plugins/org.eclipse.elk.alg.layered/src/org/eclipse/elk/alg/layered/p1cycles/DepthFirstCycleBreaker.java
//
// The ELK code is based on the following resource:
// "Emden R. Gansner, Eleftherios Koutsofios, Stephen C. North, Kiem-Phong Vo, A technique for
// drawing directed graphs. Software Engineering 19(3), pp. 214-230, 1993."
func (p *depthFirst) Process(g *graph.Graph) {
	nodeCount := len(g.Nodes)
	p.visited = make(graph.NodeSet)
	p.active = make(graph.NodeSet)

	// get list of source nodes (nodes with no incoming edge)
	sources := g.Sources()

	for _, node := range sources {
		p.visit(node)
	}

	for i := 0; i < nodeCount; i++ {
		node := g.Nodes[i]
		if !p.visited[node] {
			p.visit(node)
		}
	}

	for _, e := range p.reversable {
		e.Reverse()
		g.IsCyclic = true
	}
}

func (p *depthFirst) visit(node *graph.Node) {
	if p.visited[node] {
		return
	}
	p.visited[node] = true
	p.active[node] = true

	for _, e := range node.Out {
		// self-loop (might abstract into method)
		if e.To == node {
			continue
		}
		// Original ELK comment:
		// If the edge connects to an active node, we have found a path from said active node back to itself since
		// active nodes are on our current path. That's a backward edge and needs to be reversed
		if p.active[e.To] {
			p.reversable = append(p.reversable, e)
		} else {
			p.visit(e.To)
		}
	}

	p.active[node] = false
}

func (p *depthFirst) Cleanup() {
	p.visited = nil
	p.active = nil
	p.reversable = nil
}

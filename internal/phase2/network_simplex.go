package phase2

import (
	"math"

	"github.com/nulab/autog/internal/graph"
)

type networkSimplexProcessor struct {
	lim graph.NodeIntMap // Gansner et al.: number from a root node in spanning tree postorder traversal
	low graph.NodeIntMap // Gansner et al.: lowest postorder traversal number among nodes reachable from the input node
}

// todo: exec alg on single connected components?

// this implements a graph node layering algorithm, based on:
//   - "Emden R. Gansner, Eleftherios Koutsofios, Stephen C. North, Kiem-Phong Vo, A technique for
//     drawing directed graphs. Software Engineering 19(3), pp. 214-230, 1993."
//     https://www.researchgate.net/publication/3187542_A_Technique_for_Drawing_Directed_Graphs
//   - ELK Java code at https://github.com/eclipse/elk/blob/master/plugins/org.eclipse.elk.alg.layered/src/org/eclipse/elk/alg/layered/p2layers/NetworkSimplexLayerer.java
func execNetworkSimplex(g *graph.DGraph, params graph.Params) {
	p := &networkSimplexProcessor{
		lim: make(graph.NodeIntMap),
		low: make(graph.NodeIntMap),
	}
	p.feasibleTree(g)

	// ELK defines the max iterations as an arbitrary user value N times a fixed factor K times the sqroot of |V|.
	// where |V| is the number of nodes in each connected component.
	// N*K in ELK defaults to 28.
	k1 := int(math.Sqrt(float64(len(g.Nodes))))
	if params.NetworkSimplexMaxIterFactor > 0 {
		k1 = params.NetworkSimplexMaxIterFactor
	}
	maxitr := int(params.NetworkSimplexThoroughness) * k1

	e := negCutValueTreeEdge(g.Edges)
	i := 0
	for e != nil {
		if i >= maxitr {
			break
		}
		f := p.minSlackNonTreeEdge(g.Edges, e)
		// todo: figure out why this could be nil
		if f == nil {
			break
		}
		p.exchange(e, f, g)
		e = negCutValueTreeEdge(g.Edges)
		i++
	}
	normalize(g)
	switch params.NetworkSimplexBalance {
	case 1:
		vbalance(g)
	case 2:
		p.hbalance(g)
	}
}

// returns the first tree edge with negative cut value;
// may return nil if there is no such edge, meaning the solution is optimal.
func negCutValueTreeEdge(edges []*graph.Edge) *graph.Edge {
	for _, e := range edges {
		if e.IsInSpanningTree && e.CutValue < 0 {
			return e
		}
	}
	return nil
}

// Replace candidates are non-tree edges that go from e's head component to its tail component (original direction).
// The first candidate with minimum slack is chosen.
func (p *networkSimplexProcessor) minSlackNonTreeEdge(edges []*graph.Edge, e *graph.Edge) *graph.Edge {
	var minSlack = math.MaxInt
	var replaceCandidate *graph.Edge
	for _, f := range edges {
		if f == e || f.IsInSpanningTree {
			continue
		}
		if p.inHeadComponent(f.From, e) && !p.inHeadComponent(f.To, e) {
			slack := slack(f)
			if slack < minSlack {
				minSlack = slack
				replaceCandidate = f
			}
		}
	}
	return replaceCandidate // nil if not otherwise assigned
}

// computes an initial feasible spanning tree; it's feasible if its edges are tight
func (p *networkSimplexProcessor) feasibleTree(g *graph.DGraph) {
	p.initLayers(g)
	for {
		treeNodes := tightTree(g.Nodes[0], graph.EdgeSet{}, graph.NodeSet{})
		if len(treeNodes) == len(g.Nodes) {
			break
		}
		e := p.incidentNonTreeEdge(treeNodes)
		// incident means that one of e's vertices belongs to the tree and one doesn't.
		// here e's slack must be >= 0: since it points to a non-tree node, if the slack
		// were 0 it would've been included in the tight tree.
		// then, the layers of tree nodes are adjusted to make e's slack equal to zero.
		// as the edge becomes tight, it will be included in the tree together with its non-tree vertex
		// at the next iteration.
		d := slack(e)
		if treeNodes[e.To] {
			d = -d
		}
		for n := range treeNodes {
			n.Layer += d
		}
	}

	p.setStreeValues(g.Nodes[0])
	p.setCutValues(g)
}

// Starting from the graph source nodes (no incoming edges), this method assigns an initial layer
// to each node n based on the maximum layer of nodes connected to them via incoming edges.
// This makes all directed edges point downward.
// Source nodes have no incoming edges, so are processed first.
func (p *networkSimplexProcessor) initLayers(g *graph.DGraph) {
	// initialize the count of incoming edges for all nodes
	unseenInEdges := make(graph.NodeIntMap)
	for _, n := range g.Nodes {
		unseenInEdges[n] = n.Indeg()
	}

	// sources have layer 0
	sources := g.Sources()

	for len(sources) > 0 {
		n := sources[0]
		sources = sources[1:]

		// given a directed edge e = (n,m)
		// the target node m is assigned the layer of the source node plus delta
		// this makes the edge tight by construction because slack(e) = 0
		for _, e := range n.Out {
			m := e.To
			m.Layer = max(m.Layer, n.Layer+e.Delta)
			unseenInEdges[m]--
			if unseenInEdges[m] == 0 {
				sources = append(sources, m)
			}
		}
	}
}

// Starting from the given node, this constructs a spanning tree with only tight edges (slack = 0).
// It returns visitedNodes which contains the nodes that belong to this spanning tree.
func tightTree(n *graph.Node, visitedEdges graph.EdgeSet, visitedNodes graph.NodeSet) graph.NodeSet {
	visitedNodes[n] = true
	n.VisitEdges(func(e *graph.Edge) {
		if !visitedEdges[e] {
			visitedEdges[e] = true
			m := e.ConnectedNode(n)
			if e.IsInSpanningTree {
				tightTree(m, visitedEdges, visitedNodes)
			} else if !visitedNodes[m] && slack(e) == 0 {
				// checking that m hasn't been seen before ensures there are no loopbacks in this spanning tree
				e.IsInSpanningTree = true
				tightTree(m, visitedEdges, visitedNodes)
			}
		}
	})
	return visitedNodes
}

// This finds a "non-tree edge incident on the tree with min amount of slack".
// Incident means that only one of the edge's vertices belongs to the spanning tree.
func (p *networkSimplexProcessor) incidentNonTreeEdge(treeNodes graph.NodeSet) *graph.Edge {
	var minSlack = math.MaxInt
	var candidate *graph.Edge
	// todo: range on map non-deterministic
	for n := range treeNodes {
		n.VisitEdges(func(e *graph.Edge) {
			if e.SelfLoops() {
				return
			}
			if e.IsInSpanningTree || treeNodes[e.ConnectedNode(n)] {
				return
			}
			slack := slack(e)
			if slack < minSlack {
				minSlack = slack
				candidate = e
			}
		})
	}
	if candidate == nil {
		panic("network simplex: did not find adjacent non-tree edge with min slack: make sure the graph is connected")
	}
	return candidate
}

func (p *networkSimplexProcessor) inHeadComponent(n *graph.Node, e *graph.Edge) bool {
	if !e.IsInSpanningTree {
		panic("network simplex: breaking tree around non-tree edge")
	}
	u, v := e.From, e.To

	// the following boolean logic follows Graphviz's paper:
	// "For example, if e = (u,v) is a tree edge and vroot is in the head component of the edge (i.e., lim(u) < lim(v)),
	// then a node w is in the tail component of e if and only if low(u) ≤ lim(w) ≤ lim(u)."
	if p.lim[u] < p.lim[v] {
		if p.low[u] <= p.lim[n] && p.lim[n] <= p.lim[u] {
			// this inequality means that n is in the subtree rooted in u;
			// because of e's direction, it also implies that n is in the subtree rooted in v
			return false
		}
		return true
	}
	// else vroot is in the tail component and v is lower than u in the DFS tree
	// if n is in a subtree rooted in v, it is also in a subtree rooted in u
	// hence it's in the head component
	return p.low[v] <= p.lim[n] && p.lim[n] <= p.lim[v]
}

func (p *networkSimplexProcessor) exchange(e, f *graph.Edge, g *graph.DGraph) {
	if !e.IsInSpanningTree {
		panic("network simplex: exchange: tree-edge not in spanning tree")
	}
	if f.IsInSpanningTree {
		panic("network simplex: exchange: non-tree-edge already in spanning tree")
	}

	d := slack(f)
	if d > 0 {
		// adjust the layer of nodes in e's tail component
		for _, n := range g.Nodes {
			if !p.inHeadComponent(n, e) {
				n.Layer -= d
			}
		}
	}

	// exchange the edges
	e.IsInSpanningTree = false
	f.IsInSpanningTree = true

	// recalculate the postorder numbers and edges' cut values
	p.setStreeValues(g.Nodes[0])
	p.setCutValues(g)
}

func (p *networkSimplexProcessor) setStreeValues(n *graph.Node) {
	clear(p.lim)
	clear(p.low)
	p.walkStreeDfs(n, graph.EdgeSet{}, 1)
}

// Visits the nodes of the spanning tree in postorder traversal, assigning increasing indices.
// Same as a topological sorting; in addition, each node is mapped to a number low(n)
// which is the lowest postorder number in the subtree rooted in n.
// The root node will have low(n) = 1 and lim(n) = |V|; leaf nodes will have lim(n) = low(n).
func (p *networkSimplexProcessor) walkStreeDfs(n *graph.Node, visited graph.EdgeSet, low int) int {
	p.low[n] = low
	lim := low
	n.VisitEdges(func(e *graph.Edge) {
		if e.IsInSpanningTree && !visited[e] {
			visited[e] = true
			lim = p.walkStreeDfs(e.ConnectedNode(n), visited, lim)
		}
	})
	p.lim[n] = lim
	return lim + 1
}

// The cut value is defined as x - y where:
//   - x = sum of the weights of all edges going from the tail to the head component, including the tree edge itself
//   - y = sum of the weights of all edges from the head to the tail component
func (p *networkSimplexProcessor) setCutValues(g *graph.DGraph) {
	// todo naive implementation, optimize
	for _, e := range g.Edges {
		if !e.IsInSpanningTree {
			continue
		}
		e.CutValue += e.Weight // e itself goes from tail to head by definition

		for _, f := range g.Edges {
			// no other tree edge connects different components, otherwise we'd have two paths to e's target
			if f.IsInSpanningTree {
				continue
			}
			if !p.inHeadComponent(f.From, e) && p.inHeadComponent(f.To, e) {
				e.CutValue += f.Weight
			} else if p.inHeadComponent(f.From, e) && !p.inHeadComponent(f.To, e) {
				e.CutValue -= f.Weight
			}
		}
	}
}

func slack(e *graph.Edge) int {
	return e.To.Layer - e.From.Layer - e.Delta
}

// shifts all layers up so that the lowest layer is 0
func normalize(g *graph.DGraph) {
	lowest := math.MaxInt
	for _, n := range g.Nodes {
		lowest = min(lowest, n.Layer)
	}
	if lowest == 0 {
		return
	}
	for _, n := range g.Nodes {
		n.Layer -= lowest
	}
}

// nodes are shifted to less crowded layers if the shift preserves feasibility (edge length >= edge delta)
func vbalance(g *graph.DGraph) {
	lsize := map[int]int{}
	lmax := 0
	for _, n := range g.Nodes {
		lsize[n.Layer]++
		lmax = max(lmax, n.Layer)
	}

	for _, n := range g.Nodes {
		if n.Indeg() == n.Outdeg() {
			low := 0
			high := lmax
			for _, e := range n.In {
				low = max(low, e.From.Layer+e.Delta)
			}
			for _, e := range n.Out {
				high = min(high, e.To.Layer-e.Delta)
			}
			newl := low

			// if the node has only flat edges, or in/out-span 1, or is source/sink with span 1, this does nothing
			// otherwise it may shift the node
			for i := low + 1; i <= high; i++ {
				if lsize[i] < lsize[newl] {
					newl = i
				}
			}
			if lsize[newl] < lsize[n.Layer] {
				lsize[n.Layer]--
				lsize[newl]++
				n.Layer = newl
			}
		}
	}
}

func (p *networkSimplexProcessor) hbalance(g *graph.DGraph) {
	for _, e := range g.Edges {
		if !e.IsInSpanningTree {
			continue
		}
		if e.CutValue == 0 {
			f := p.minSlackNonTreeEdge(g.Edges, e)
			if f == nil {
				continue
			}
			d := slack(f)
			if d < 1 {
				continue
			}
			if p.lim[e.From] < p.lim[e.To] {
				p.adjustLayers(e.From, d)
			} else {
				p.adjustLayers(e.To, -d)
			}
		}
	}
}

func (p *networkSimplexProcessor) adjustLayers(n *graph.Node, delta int) {
	n.Layer -= delta
	for _, e := range n.Out {
		if !e.IsInSpanningTree {
			continue
		}
		if !(p.lim[n] < p.lim[e.ConnectedNode(n)]) {
			p.adjustLayers(e.To, delta)
		}
	}
	for _, e := range n.In {
		if !e.IsInSpanningTree {
			continue
		}
		if !(p.lim[n] < p.lim[e.ConnectedNode(n)]) {
			p.adjustLayers(e.From, delta)
		}
	}
}

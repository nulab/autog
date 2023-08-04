package phase2

import (
	"math"

	"github.com/nulab/autog/graph"
)

const (
	thoroughness = 28
)

type networkSimplexProcessor struct {
	poIndex int              // node post-order traversal index
	lim     graph.NodeIntMap // Gansner et al.: number from a root node in spanning tree postorder traversal
	low     graph.NodeIntMap // Gansner et al.: lowest postorder traversal number among nodes reachable from the input node
}

// todo: exec alg on single connected components?

// this implements a graph node layering algorithm, based on:
//   - "Emden R. Gansner, Eleftherios Koutsofios, Stephen C. North, Kiem-Phong Vo, A technique for
//     drawing directed graphs. Software Engineering 19(3), pp. 214-230, 1993."
//     https://www.researchgate.net/publication/3187542_A_Technique_for_Drawing_Directed_Graphs
//   - ELK Java code at https://github.com/eclipse/elk/blob/master/plugins/org.eclipse.elk.alg.layered/src/org/eclipse/elk/alg/layered/p2layers/NetworkSimplexLayerer.java
func execNetworkSimplex(g *graph.DGraph) {
	p := &networkSimplexProcessor{
		lim: make(graph.NodeIntMap),
		low: make(graph.NodeIntMap),
	}
	p.feasibleTree(g)

	// ELK defines the max iterations as an arbitrary user value N times a fixed factor K times the sqroot of |V|.
	// where |V| is the number of nodes in each connected component.
	// N*K in ELK defaults to 28.
	maxitr := thoroughness * int(math.Sqrt(float64(len(g.Nodes))))

	e := negCutValueTreeEdge(g.Edges)
	i := 0
	for e != nil {
		if i >= maxitr {
			break
		}
		f := p.minSlackNonTreeEdge(g.Edges, e)
		p.exchange(e, f, g)
		e = negCutValueTreeEdge(g.Edges)
		i++
	}
	normalize(g)
	balance(g)
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
	p.poIndex = 1
	p.postOrderTraversal(g.Nodes[0], graph.EdgeSet{})
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
	for n := range treeNodes {
		n.VisitEdges(func(e *graph.Edge) {
			if e.ConnectedNode(n) == n {
				return // avoid self-loops
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
		panic("network simplex: did not find adjacent non-tree edge with min slack")
	}
	return candidate
}

func (p *networkSimplexProcessor) inHeadComponent(n *graph.Node, e *graph.Edge) bool {
	if !e.IsInSpanningTree {
		panic("network simplex: breaking tree around non-tree edge")
	}
	u, _ := p.postorderOf(e)
	// Note that lim(n) < lim(u) means that n was visited before u in postorder traversal
	// The first condition alone is not enough because n could belong to a tree branch to the left of e.
	inTail := p.low[u] <= p.lim[n] && p.lim[n] <= p.lim[u]
	return !inTail
}

func (p *networkSimplexProcessor) exchange(e, f *graph.Edge, g *graph.DGraph) {
	if !e.IsInSpanningTree {
		panic("network simplex: exchange: tree-edge not in spanning tree")
	}
	if f.IsInSpanningTree {
		panic("network simplex: exchange: non-tree-edge already in spanning tree")
	}

	ftail, fhead := p.postorderOf(f)

	d := ftail.Layer - fhead.Layer - e.Delta
	if !p.inHeadComponent(ftail, e) {
		d *= -1
	}

	// adjust the layer of nodes in e's tail component
	for _, n := range g.Nodes {
		if !p.inHeadComponent(n, e) {
			n.Layer += d
		}
	}

	// exchange the edges
	e.IsInSpanningTree = false
	f.IsInSpanningTree = true

	// recalculate the postorder numbers and edges' cut values
	p.poIndex = 1
	p.postOrderTraversal(g.Nodes[0], graph.EdgeSet{})
	p.setCutValues(g)
}

// Visits the nodes of the spanning tree in postorder traversal, assigning increasing indices.
// Same as a topological sorting; in addition, each node is mapped to a number low(n)
// which is the lowest postorder number in the subtree rooted in n.
// The root node will have low(n) = 1 and lim(n) = |V|; leaf nodes will have lim(n) = low(n).
func (p *networkSimplexProcessor) postOrderTraversal(n *graph.Node, visited graph.EdgeSet) int {
	if len(visited) == 0 && p.poIndex != 1 {
		panic("network simplex: must initialize postorder ordinal number")
	}
	lowest := math.MaxInt
	for _, e := range n.Edges() {
		if e.IsInSpanningTree && !visited[e] {
			visited[e] = true
			lowest = min(lowest, p.postOrderTraversal(e.ConnectedNode(n), visited))
		}
	}
	p.lim[n] = p.poIndex
	p.low[n] = min(lowest, p.poIndex)
	p.poIndex++
	return p.low[n]
}

// Gansner et al.:
// "For each tree edge, [the cut value] is computed by marking the nodes as belonging to the head or tail component,
// and then performing the sum of the signed weights of all edges whose [source] and [target] nodes are in different components,
// the sign [of the weight] being negative for those edges going from the head to the tail component."
func (p *networkSimplexProcessor) setCutValues(g *graph.DGraph) {
	// todo naive implementation, optimize

	// The cut value is:
	// sum of the weights of all edges going from the tail to the head component, including the tree edge itself
	// minus
	// sum of the weights of all edges from the head to the tail component.
	for _, e := range g.Edges {
		if !e.IsInSpanningTree {
			continue
		}
		th := 0
		ht := 0

		th += e.Weight // e goes from tail to head by construction
		// no other tree edge connects different components

		for _, f := range g.Edges {
			if f.IsInSpanningTree {
				continue
			}
			if p.lim[f.To] <= p.lim[e.From] /* f.To is in tail component */ {
				if p.lim[f.From] >= p.lim[e.To] /* f.From is in head component */ {
					ht -= e.Weight
				}

			} else /* f.To is in head component */ {
				if p.lim[f.From] <= p.lim[e.From] /* f.From is in tail component */ {
					th += e.Weight
				}
			}
		}
		e.CutValue = th - ht
	}
}

func slack(e *graph.Edge) int {
	return e.To.Layer - e.From.Layer - e.Delta
}

// The definition of head and tail in Gansner et al.'s paper are relative to the root of the postorder traversal:
// after "deleting" the edge e, the head component is the one that contains the root of the tree.
// Therefore, the direction of the *graph.Edge (From->To) in the source graph doesn't necessarily
// reflect the direction of the edge in postorder traversal.
// The postorder direction depends on direction of the inequality between lim(a) and lim(b).
// Gansner et al. p.219: "For example, if e = (u,v) is a tree edge and vroot is in the head component of the edge (i.e., lim(u) < lim(v)),
// then a node w is in the tail component of e if and only if low(u) <= lim(w) <= lim(u)".
//
// Also, the quote above intuitively means that w is in the tail component if it has u as ancestor in postorder direction.
//
// Now, a node u is in the tail component of e if lim(u) < lim(v)
// so among the *graph.Edge x and y poles, u is whichever of the two has a lower lim()
func (p *networkSimplexProcessor) postorderOf(e *graph.Edge) (tailNode, headNode *graph.Node) {
	x, y := e.From, e.To
	if p.lim[x] < p.lim[y] {
		return x, y
	}
	return y, x
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
func balance(g *graph.DGraph) {
	lsize := map[int]int{}
	for _, n := range g.Nodes {
		lsize[n.Layer]++
	}
	for _, n := range g.Nodes {
		if n.Indeg() == n.Outdeg() {
			l := n.Layer

			span := feasibleSpan(n)
			// if the node has only flat edges, or in/out-span 1, or is source/sink with span 1, this does nothing
			// otherwise it may shift the node
			for i := l - span[0] + 1; i < l+span[1]; i++ {
				if lsize[i] < lsize[l] {
					l = i
				}
			}
			// node could've moved back to the original layer
			if lsize[l] < lsize[n.Layer] {
				lsize[n.Layer]--
				lsize[l]++
				n.Layer = l
			}

		}
	}
}

func feasibleSpan(n *graph.Node) (span [2]int) {
	minInSpan := math.MaxInt
	minOutSpan := math.MaxInt

	n.VisitEdges(func(e *graph.Edge) {
		edgespan := e.To.Layer - e.From.Layer
		switch {
		case e.To == n && edgespan < minInSpan:
			minInSpan = edgespan

		case e.From == n && edgespan < minOutSpan:
			minOutSpan = edgespan
		}
	})
	if minInSpan == math.MaxInt {
		minInSpan = -1
	}
	if minOutSpan == math.MaxInt {
		minOutSpan = -1
	}

	span[0] = minInSpan
	span[1] = minOutSpan
	return
}

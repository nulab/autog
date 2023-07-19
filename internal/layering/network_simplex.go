package layering

import (
	"fmt"
	"math"

	"github.com/vibridi/autog/internal/graph"
)

type nodemap map[*graph.Node]int // todo: tidy up, also used in greedy cyclebreaker
type edgemap map[*graph.Edge]int

const (
	delta  = 1
	maxitr = 25
)

type networkSimplex struct {
	// visited   graph.NodeSet
	postOrder int
	lim       nodemap // Emden et al.: number from a root node in spanning tree postorder traversal
	low       nodemap // Emden et al.: lowest postorder traversal number among nodes reachable from the input node
}

var NetworkSimplex = &networkSimplex{}

// Process implements a graph node layering algorithm, based on:
//   - "Emden R. Gansner, Eleftherios Koutsofios, Stephen C. North, Kiem-Phong Vo, A technique for
//     drawing directed graphs. Software Engineering 19(3), pp. 214-230, 1993."
//     https://www.researchgate.net/publication/3187542_A_Technique_for_Drawing_Directed_Graphs
//   - ELK Java code at https://github.com/eclipse/elk/blob/master/plugins/org.eclipse.elk.alg.layered/src/org/eclipse/elk/alg/layered/p2layers/NetworkSimplexLayerer.java
func (p *networkSimplex) Process(g *graph.DGraph) {
	p.lim = nodemap{}
	p.low = nodemap{}

	p.feasibleTree(g)
	fmt.Println(graph.EdgeList(g.Edges))

	e := negCutValueTreeEdge(g.Edges)
	for i := 0; e != nil && i < maxitr; i++ {
		f := p.minSlackNonTreeEdge(g.Edges, e)
		p.exchange(e, f, g)
		e = negCutValueTreeEdge(g.Edges)
	}
	// normalize()
	// balance()
}

func (p *networkSimplex) feasibleTree(g *graph.DGraph) {
	fmt.Println("init layers")
	p.initLayers(g)
	// todo improve readability
	i := 0
	for {
		i++
		treeNodes := graph.NodeSet{}
		numNodes := tightTree(g.Nodes[0], graph.EdgeSet{}, treeNodes)
		if numNodes >= len(g.Nodes) { // todo maybe just use equal
			break
		}
		if i >= maxitr*2 {
			fmt.Println("feasibleTree max iter reached")
			break
		}
		// todo this loop-with-break construct could be replaced with initializing the treeNodes map in a
		// 	for loop clause and then use clear at the end of the loop, needs Go 1.21
		// This finds an edge to a nontree node that is adjacent
		// to the tree, and adjusts the ranks of the tree nodes to make
		// this edge tight. As the edge was picked to have minimal
		// slack, the resulting ranking is still feasible. Thus, on every
		// iteration, the maximal tight tree gains at least one node, and the
		// algorithm eventually terminates with a feasible spanning tree.
		// This technique is essentially the one described by Sugiyama
		// et a1 [5].
		e := p.adjacentNonTreeEdge(treeNodes)
		d := slack(e)
		if treeNodes[e.To] {
			d = -d
		}
		for n := range treeNodes {
			n.Layer += d
		}

	}
	p.postOrderTraversal(g.Nodes[0], graph.EdgeSet{})
	p.setCutValues(g)
}

// returns the first tree edge with negative cut value;
// may return nil if there is no such edge, meaning the solution is optimal.
func negCutValueTreeEdge(t graph.EdgeList) *graph.Edge {
	for _, e := range t {
		if e.IsInSpanningTree && e.CutValue < 0 {
			return e
		}
	}
	return nil
}

// Finds a non-tree edge to replace e.
// This is done by breaking the tree into head and tail components around the edge e,
// then picking the non-tree edge with minimum slack that goes from head to tail.
// Note that the argument t is the spanning tree but contains also non-tree edges marked as such.
func (p *networkSimplex) minSlackNonTreeEdge(t graph.EdgeList, e *graph.Edge) *graph.Edge {
	var minSlack = math.MaxInt
	var replaceCandidate *graph.Edge
	for _, f := range t {
		if f == e || f.IsInSpanningTree {
			continue
		}
		tailNode, headNode := p.postorderOf(f)
		if p.inHeadComponent(headNode, e) && !p.inHeadComponent(tailNode, e) {
			slack := slack(f)
			if slack < minSlack {
				minSlack = slack
				replaceCandidate = f
			}
		}
	}
	return replaceCandidate // nil if not otherwise assigned
}

// "non-tree edge incident on the tree with min amount of slack"
// incident means that the tree node could be either the edge's source or target
func (p *networkSimplex) adjacentNonTreeEdge(treeNodes graph.NodeSet) *graph.Edge {
	// This finds an edge to a nontree node that is adjacent
	// to the tree
	minSlack := math.MaxInt
	var candidate *graph.Edge
	for n := range treeNodes { // todo iterate over map, non-deterministic
		for itr := n.EdgeIter(); itr.HasNext(); {
			e := itr.Next()
			if e.ConnectedNode(n) == n {
				continue // avoid self-loops
			}
			slack := slack(e)
			if slack < minSlack {
				minSlack = slack
				candidate = e
			}
		}
	}
	if candidate == nil {
		panic("network simplex: did not find adjacent non-tree edge with min slack")
	}
	return candidate
}

func (p *networkSimplex) inHeadComponent(n *graph.Node, e *graph.Edge) bool {
	if !e.IsInSpanningTree {
		panic("network simplex: breaking tree around non-tree edge")
	}
	u, _ := p.postorderOf(e)
	// Note that lim(n) < lim(u) means that n was visited before u in postorder traversal
	// The first condition alone is not enough because n could belong to a tree branch to the left of e.
	inTail := p.low[u] <= p.lim[n] && p.lim[n] <= p.lim[u]
	return !inTail
}

// The edges are exchanged, updating the tree and its cut
// values.
func (p *networkSimplex) exchange(e, f *graph.Edge, g *graph.DGraph) {
	if !e.IsInSpanningTree {
		panic("exchange: tree-edge not in spanning tree")
	}
	if f.IsInSpanningTree {
		panic("exchange: non-tree-edge already in spanning tree")
	}

	ftail, fhead := p.postorderOf(f)

	d := ftail.Layer - fhead.Layer - delta
	if !p.inHeadComponent(ftail, e) {
		d *= -1
	}

	// adjust the layer of nodes in e's tail component
	for _, n := range g.Nodes {
		if !p.inHeadComponent(n, e) {
			n.Layer += d
		}
	}

	e.IsInSpanningTree = false
	f.IsInSpanningTree = true

	p.postOrder = 1                                   // todo: why reset this here?
	p.postOrderTraversal(g.Nodes[0], graph.EdgeSet{}) // todo: node, elk uses graph.nodes.iterator().next()
	p.setCutValues(g)
}

// Visits the nodes of the spanning tree in postorder traversal, assigning increasing postorder numbers
// starting from the farthest child from the root n and proceeding in DFS order.
// low(n) is the lowest postorder number in the subtree rooted in n.
// The root node will have low(n) = 1 and lim(n) = |V|; leaf nodes will have lim(n) = low(n).
func (p *networkSimplex) postOrderTraversal(n *graph.Node, visited graph.EdgeSet) int {
	lowest := math.MaxInt
	for _, e := range n.Edges() {
		if e.IsInSpanningTree && !visited[e] {
			visited[e] = true
			lowest = min(lowest, p.postOrderTraversal(e.ConnectedNode(n), visited))
		}
	}
	p.lim[n] = p.postOrder
	p.low[n] = min(lowest, p.postOrder)
	p.postOrder++
	return p.low[n]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Nodes having equal in- and out-edge weights [NDR: degrees]
// and multiple feasible ranks are moved to a feasible rank with the fewest
// nodes.
// [NDR: feasible rank is such that l(e) >= d(e)]
// The purpose is to reduce crowding and improve the
// aspect ratio of the drawing, following principle A4. The adjustment does not change the cost of the rank assignment. Nodes
// are adjusted in a greedy fashion, which works sufficiently well.
func balance() {

}

// Starting from the graph source nodes (no incoming edges), this method assigns an initial layer
// to each node n based on the maximum layer of nodes connected to them via incoming edges.
// Source nodes have no incoming edges, so are processed first.
func (p *networkSimplex) initLayers(g *graph.DGraph) {
	// initialize the count of incoming edges for all nodes
	unseenInEdges := nodemap{}
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
			m.Layer = max(m.Layer, n.Layer+delta)
			unseenInEdges[m]--
			if unseenInEdges[m] == 0 {
				sources = append(sources, m)
			}
		}
	}
}

// Function tight-tree finds a maximal tree of tight
// edges containing some fixed node and returns the number of
// nodes in the tree. Note that such a maximal tree is just a
// spanning tree for the subgraph induced by all nodes reachable
// from the fixed node in the underlying undirected graph using
// only tight edges. In particular, all such trees have the same
// number of nodes.
func tightTree(n *graph.Node, visitedEdges graph.EdgeSet, visitedNodes graph.NodeSet) int {
	nodeCount := 1
	visitedNodes[n] = true
	for itr := n.EdgeIter(); itr.HasNext(); {
		e := itr.Next()

		if !visitedEdges[e] {
			visitedEdges[e] = true
			m := e.ConnectedNode(n)
			if e.IsInSpanningTree {
				nodeCount += tightTree(m, visitedEdges, visitedNodes)
			} else if !visitedNodes[m] && slack(e) == 0 {
				// checking that m hasn't been seen before ensures there are no loopbacks in this spanning tree
				e.IsInSpanningTree = true
				nodeCount += tightTree(m, visitedEdges, visitedNodes)
			}
		}

	}
	return nodeCount
}

// The init-cutvalues function computes the cut values of the tree edges. For each tree edge, this is computed by
// marking the nodes as belonging to the head or tail component,

// head are all nodes on the side of the edge's target node (always true by construction)
// tail are all nodes on the side of the edge's source node

// The cut value is the
// sum of the weights of all edges going from the tail to the head component, including the tree edge itself
// minus
// sum of the weights of all edges from the head to the tail component.

// and then performing the sum of the signed weights of all
// edges whose head and tail are in different components, the
// sign being negative for those edges going from the head to
// the tail component.

// Emden et al.:
// "For each tree edge, [the cut value] is computed by marking the nodes as belonging to the head or tail component,
// and then performing the sum of the signed weights of all edges whose [source] and [target] nodes are in different components,
// the sign [of the weight] being negative for those edges going from the head to the tail component."
func (p *networkSimplex) setCutValues(g *graph.DGraph) {
	// todo naive implementation, optimize

	// The cut value is the
	// sum of the weights of all edges going from the tail to the head component, including the tree edge itself
	// minus
	// sum of the weights of all edges from the head to the tail component.
	for _, e := range g.Edges {
		if !e.IsInSpanningTree {
			continue
		}
		th := 0
		ht := 0

		th += delta // e goes from tail to head by construction
		// no other tree edge connects different components

		for _, f := range g.Edges {
			if f.IsInSpanningTree {
				continue
			}
			if p.lim[f.To] <= p.lim[e.From] /* f.To is in tail component */ {
				if p.lim[f.From] >= p.lim[e.To] /* f.From is in head component */ {
					ht -= delta
				}

			} else /* f.To is in head component */ {
				if p.lim[f.From] <= p.lim[e.From] /* f.From is in tail component */ {
					th += delta
				}
			}
		}
		e.CutValue = th - ht
	}
}

func slack(e *graph.Edge) int {
	return e.To.Layer - e.From.Layer - delta
}

// The definition of head and tail in Emden et al.'s paper are relative to the root of the postorder traversal:
// after "deleting" the edge e, the head component is the one that contains the root of the tree.
// Therefore, the direction of the *graph.Edge (From->To) in the source graph doesn't necessarily
// reflect the direction of the edge in postorder traversal.
// The postorder direction depends on direction of the inequality between lim(a) and lim(b).
// Emden et al. p.219: "For example, if e = (u,v) is a tree edge and vroot is in the head component of the edge (i.e., lim(u) < lim(v)),
// then a node w is in the tail component of e if and only if low(u) <= lim(w) <= lim(u)".
//
// Also, the quote above intuitively means that w is in the tail component if it has u as ancestor in postorder direction.
//
// Now, a node u is in the tail component of e if lim(u) < lim(v)
// so among the *graph.Edge x and y poles, u is whichever of the two has a lower lim()
func (p *networkSimplex) postorderOf(e *graph.Edge) (tailNode, headNode *graph.Node) {
	x, y := e.From, e.To
	if p.lim[x] < p.lim[y] {
		return x, y
	}
	return y, x
}

func tight(e *graph.Edge) bool {
	return slack(e) == 0
}

func (p *networkSimplex) Cleanup() {
	p.lim = nil
	p.low = nil
	p.postOrder = 0

}

package phase1

import (
	"math"
	"math/rand"
	"time"

	"github.com/nulab/autog/graph"
)

type greedyProcessor struct {
	rnd     *rand.Rand
	arcdiag graph.NodeMap
	outdeg  graph.NodeMap
	indeg   graph.NodeMap
}

// Process implements a greedy cycle breaker.
// This is a port of ELK Java code, slightly modified to adapt to Go data structures and coding practices.
// https://github.com/eclipse/elk/blob/master/plugins/org.eclipse.elk.alg.layered/src/org/eclipse/elk/alg/layered/p1cycles/GreedyCycleBreaker.java
//
// The ELK code is based on the following resources:
//   - Peter Eades, Xuemin Lin, W. F. Smyth: A fast and effective heuristic for the feedback arc set problem
//     http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.47.7745
//   - Giuseppe di Battista, Peter Eades, Roberto Tamassia, Ioannis G. Tollis,
//     Graph Drawing: Algorithms for the Visualization of Graphs, Prentice Hall, New Jersey, 1999 (Section 9.4).
//
// The algorithm arranges the nodes of G in an arc diagram, with source nodes to the right and sink nodes to the left.
// Then it reverses edges that point right.
func execGreedy(g *graph.DGraph) {
	p := greedyProcessor{
		rnd:     rand.New(rand.NewSource(time.Now().UnixNano())),
		arcdiag: make(graph.NodeMap),
		outdeg:  make(graph.NodeMap),
		indeg:   make(graph.NodeMap),
	}

	nodeCount := len(g.Nodes)

	sources := g.Sources()
	sinks := g.Sinks()

	// here ELK accounts for edge priority: particular edges that the user doesn't want to reverse
	// can be assigned a non-zero priority; this will artificially increase the node's in-/out-degrees.

	// numerical ranks for sinks and sources
	// nodes will be arranged in the arc diagram in decreasing order from left to right based on their indegree
	nextRight := -1
	nextLeft := 1

	var maxOutflowNodes []*graph.Node

	// arrange nodes in the arc diagram
	for i := nodeCount; i > 0; {
		// sinks are put to the right of the arc diagram --> assign negative rank, which is later shifted to positive
		for len(sinks) > 0 {
			sink := sinks[0]
			sinks = sinks[1:]

			p.arcdiag[sink] = nextRight
			nextRight--
			p.updateNeighbors(sink, &sources, &sinks)
			i--
		}

		// sources are put to the left of the arc diagram
		for len(sources) > 0 {
			source := sources[0]
			sources = sources[1:]

			p.arcdiag[source] = nextLeft
			nextLeft++
			p.updateNeighbors(source, &sources, &sinks)
			i--
		}

		// while there are unprocessed nodes left that are neither sinks nor sources...
		for i > 0 {
			maxOutflow := math.MinInt

			for _, n := range g.Nodes {
				if _, ok := p.arcdiag[n]; ok {
					// already processed
					continue
				}
				outflow := p.outdeg[n] - p.indeg[n]
				if outflow >= maxOutflow {
					if outflow > maxOutflow {
						maxOutflowNodes = nil
						maxOutflow = outflow
					}
					maxOutflowNodes = append(maxOutflowNodes, n)
				}
			}
			if !(maxOutflow > math.MinInt) {
				panic("expected maxOutflow strictly greater than MinInt")
			}

			// randomly select a node from the ones with maximal outflow and put it left
			n := p.pickRandom(maxOutflowNodes)
			p.arcdiag[n] = nextLeft
			nextLeft++
			p.updateNeighbors(n, &sources, &sinks)
			i--
		}
	}

	// shift negative ranks to positive; this applies to sinks of the graph
	shift := nodeCount + 1
	for _, n := range g.Nodes {
		if p.arcdiag[n] < 0 /* sink node */ {
			p.arcdiag[n] += shift
		}
	}

	// reverse edges that point right
	for _, n := range g.Nodes {
		for _, e := range n.Out {
			if p.arcdiag[n] > p.arcdiag[e.To] {
				e.Reverse()
			}
		}
	}
}

func (p *greedyProcessor) pickRandom(nodes []*graph.Node) *graph.Node {
	return nodes[len(nodes)/2] // todo: deterministic for debugging
	// return nodes[p.rnd.Intn(len(nodes))]
}

// Updates indegree and outdegree values of the neighbors of the given node,
// simulating its removal from the graph. the sources and sinks lists are also updated.
func (p *greedyProcessor) updateNeighbors(n *graph.Node, sources, sinks *[]*graph.Node) {
	for _, e := range n.In {
		if e.SelfLoops() {
			continue
		}
		src := e.From
		if _, ok := p.arcdiag[src]; ok {
			// already processed
			continue
		}
		p.outdeg[src]--
		if p.outdeg[src] <= 0 && p.indeg[src] > 0 {
			*sinks = append(*sinks, src)
		}
	}
	for _, e := range n.Out {
		if e.SelfLoops() {
			continue
		}
		tgt := e.To
		if _, ok := p.arcdiag[tgt]; ok {
			// already processed
			continue
		}
		p.indeg[tgt]--
		if p.indeg[tgt] <= 0 && p.outdeg[tgt] > 0 {
			*sources = append(*sources, tgt)
		}
	}
}

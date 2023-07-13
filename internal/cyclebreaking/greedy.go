package cyclebreaking

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/vibridi/autog/internal/graph"
	"golang.org/x/exp/slices"
)

// todo: refactor somehow
type nodemap map[*graph.Node]int

func (m nodemap) String() string {
	type pair struct {
		n *graph.Node
		i int
	}
	var kvPairs []pair
	for k, v := range m {
		kvPairs = append(kvPairs, pair{k, v})
	}
	slices.SortFunc(kvPairs, func(a, b pair) bool {
		return a.i > b.i
	})
	bld := strings.Builder{}
	for _, p := range kvPairs {
		bld.WriteRune('[')
		bld.WriteString(p.n.ID)
		bld.WriteRune(':')
		bld.WriteString(strconv.Itoa(p.i))
		bld.WriteRune(']')
		bld.WriteRune(' ')
	}
	return bld.String()
}

type greedy struct {
	rnd     *rand.Rand
	arcdiag nodemap
	outdeg  nodemap
	indeg   nodemap
}

var Greedy = &greedy{}

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
// Then it reverses edges that point left.
func (p *greedy) Process(g *graph.DGraph) {
	p.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	p.arcdiag = nodemap{}
	p.outdeg = nodemap{}
	p.indeg = nodemap{}

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

func (p *greedy) pickRandom(nodes []*graph.Node) *graph.Node {
	return nodes[p.rnd.Intn(len(nodes))]
}

// Updates indegree and outdegree values of the neighbors of the given node,
// simulating its removal from the graph. the sources and sinks lists are also updated.
func (p *greedy) updateNeighbors(n *graph.Node, sources, sinks *[]*graph.Node) {
	for _, e := range n.In {
		src := e.From
		if src == n /* self-loop */ {
			continue
		}
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
		tgt := e.To
		if tgt == n /* self-loop */ {
			continue
		}
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

func (p *greedy) Cleanup() {
	p.rnd = nil
	p.arcdiag = nil
	p.outdeg = nil
	p.indeg = nil
}

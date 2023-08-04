package phase3

import (
	"sort"
	"strconv"

	"github.com/nulab/autog/graph"
)

type initDirection uint8

const (
	initDirectionTop    initDirection = 0
	initDirectionBottom initDirection = 1
)

type graphvizDotProcessor struct {
	positions      graph.NodeIntMap
	flipEqual      bool
	transposeEqual bool
	fixedPositions fixedPositions
}

// Ordering algorithm used in Graphviz Dot and described in:
//   - "Emden R. Gansner, Eleftherios Koutsofios, Stephen C. North, Kiem-Phong Vo, A technique for
//     drawing directed graphs. Software Engineering 19(3), pp. 214-230, 1993."
//     https://www.researchgate.net/publication/3187542_A_Technique_for_Drawing_Directed_Graphs
//
// Note that ELK's implementation is based on the original algorithm proposed by Sugiyama et al. instead of Graphviz.
func execGraphvizDot(g *graph.DGraph, params graph.Params) {
	if len(g.Layers) == 1 {
		// no crossings to reduce
		return
	}

	// insert virtual nodes so that edges with length >1 have length 1
	breakLongEdges(g)

	p3monitor := phase3monitor{"graphvizdot", params.Monitor}

	maxiter := params.GraphvizDotMaxIter
	fixedPositions := initFixedPositions(g.Edges)

	bestx_top, bestpos_top := graphvizRun(g, graphvizRunParams{maxiter, fixedPositions, initDirectionTop})
	bestx_btm, bestpos_btm := graphvizRun(g, graphvizRunParams{maxiter, fixedPositions, initDirectionBottom})

	var (
		bestx                  = 0
		bestp graph.NodeIntMap = nil
	)
	if bestx_top < bestx_btm {
		bestx, bestp = bestx_top, bestpos_top
	} else {
		bestx, bestp = bestx_btm, bestpos_btm
	}

	p3monitor.Send("crossings", bestx)

	// reset the best node positions using the saved bestp
	for _, n := range g.Nodes {
		n.LayerPos = bestp[n]
	}
	for _, l := range g.Layers {
		sort.Slice(l.Nodes, func(i, j int) bool {
			return l.Nodes[i].LayerPos < l.Nodes[j].LayerPos
		})
	}
}

func breakLongEdges(g *graph.DGraph) {
	v := 1
	i := 0
loop:
	for i < len(g.Edges) {
		e := g.Edges[i]
		i++
		if e.To.Layer-e.From.Layer > 1 {
			from, to := e.From, e.To
			// create virtual node
			virtualNode := &graph.Node{
				ID:        "V" + strconv.Itoa(v),
				Layer:     from.Layer + 1,
				IsVirtual: true,
				Size:      graph.Size{H: 100.0, W: 100.0}, // todo: eventually this doesn't belong here
			}
			v++
			// set e's target to the virtual node
			e.To = virtualNode
			// add e to virtual node incoming edges
			virtualNode.In = append(virtualNode.In, e)
			// create new edge from virtual to e's former target
			f := graph.NewEdge(virtualNode, to, 1)
			f.IsReversed = e.IsReversed
			// add f to virtual node outgoing edges
			virtualNode.Out = []*graph.Edge{f}
			// replace e with f in e's former target incoming edges
			for _, in := range to.In {
				if in == e {
					e = f
				}
			}
			// update the graph's node and edge lists
			g.Edges = append(g.Edges, f)
			g.Nodes = append(g.Nodes, virtualNode)
			g.Layers[virtualNode.Layer].Nodes = append(g.Layers[virtualNode.Layer].Nodes, virtualNode)
			// restart loop
			goto loop
		}
	}
}

type graphvizRunParams struct {
	maxiter        int
	fixedPositions fixedPositions
	dir            initDirection
}

// node order is maintained in three different places:
//   - in g.Layers.Nodes, which is a slice
//   - in each node.LayerPos field
//   - in p.positions
//
// at each iteration, this algorithm will update the node positions in all three places
// a copy of the best p.positions is kept and at the end it is propagated to g.Layers and node.LayerPos
func graphvizRun(g *graph.DGraph, params graphvizRunParams) (int, graph.NodeIntMap) {
	p := &graphvizDotProcessor{
		positions:      graph.NodeIntMap{},
		fixedPositions: params.fixedPositions,
	}
	switch params.dir {
	case initDirectionTop:
		p.initTop(g)
	case initDirectionBottom:
		p.initBottom(g)
	}
	layers := g.Layers // shallow copy

	// propagate initial order to g.Layers.Nodes slice order
	for _, layer := range layers {
		sort.Slice(layer.Nodes, func(i, j int) bool {
			a, b := layer.Nodes[i], layer.Nodes[j]
			if a.Layer != b.Layer {
				panic("same-layer nodes have different layers")
			}
			return p.positions[a] < p.positions[b]
		})
	}

	bestx := crossings(layers)
	bestp := p.positions.Clone()

	for i := 0; i < params.maxiter; i++ {
		// Depending on the parity of the current iteration
		// number, the ranks are traversed from top to bottom or from bottom to top.
		if i%2 == 0 {
			p.wmedianTopBottom(layers)
		} else {
			p.wmedianBottomTop(layers)
			p.flipEqual = !p.flipEqual // switch after every two iterations
		}
		p.transpose(layers)
		p.transposeEqual = !p.transposeEqual // switch after every two iterations

		// todo: adaptive strategy to keep iterating in case of sufficiently large improvement
		if x := crossings(layers); x < bestx {
			bestx = x
			bestp = p.positions.Clone()
		}
		if bestx == 0 {
			break
		}
	}
	return bestx, bestp
}

func (p *graphvizDotProcessor) initTop(g *graph.DGraph) {
	// initialize positions
	visited := graph.NodeSet{}
	indices := map[int]int{}
	for _, n := range g.Layers[0].Nodes {
		p.initPositionsFromTop(n, visited, indices)
	}
	for _, n := range g.Nodes {
		p.initPositionsFromTop(n, visited, indices)
	}
}

func (p *graphvizDotProcessor) initBottom(g *graph.DGraph) {
	// initialize positions
	visited := graph.NodeSet{}
	indices := map[int]int{}
	for _, n := range g.Layers[len(g.Layers)-1].Nodes {
		p.initPositionsFromBottom(n, visited, indices)
	}
	for _, n := range g.Nodes {
		p.initPositionsFromBottom(n, visited, indices)
	}
}

func (p *graphvizDotProcessor) initPositionsFromTop(n *graph.Node, visited graph.NodeSet, indices map[int]int) {
	if visited[n] {
		return
	}
	visited[n] = true
	p.setPos(n, indices[n.Layer])
	indices[n.Layer]++
	p.initPositionsFlatEdges(n, visited, indices)
	for _, e := range n.Out {
		p.initPositionsFromTop(e.To, visited, indices)
	}
}

func (p *graphvizDotProcessor) initPositionsFromBottom(n *graph.Node, visited graph.NodeSet, indices map[int]int) {
	if visited[n] {
		return
	}
	visited[n] = true
	p.setPos(n, indices[n.Layer])
	indices[n.Layer]++
	p.initPositionsFlatEdges(n, visited, indices)
	for _, e := range n.In {
		p.initPositionsFromBottom(e.From, visited, indices)
	}
}

func (p *graphvizDotProcessor) initPositionsFlatEdges(n *graph.Node, visited graph.NodeSet, indices map[int]int) {
	h, i := p.fixedPositions.head(n)
	if i > 0 {
		for h != nil && h != n {
			visited[h] = true
			p.setPos(h, indices[h.Layer])
			indices[h.Layer]++
			h = p.fixedPositions.mustBefore[h]
		}
	}
}

// The weighted median routine assigns an order to each vertex in layer L(i) based on the current order
// of adjacent nodes in the next rank. Next is L(i)-1 in top-bottom sweep, or L(i)+1 in bottom-top sweep.
// Nodes with no adjacent nodes in the next layer are kept in place.
func (p *graphvizDotProcessor) wmedianTopBottom(layers map[int]*graph.Layer) {
	medians := graph.NodeFloatMap{}
	for r := 1; r < len(layers); r++ {
		for _, v := range layers[r].Nodes {
			medians[v] = medianOf(p.adjacentNodesPositions(v, v.In, r-1))
		}
		p.sortLayer(layers[r].Nodes, medians)
	}
}

func (p *graphvizDotProcessor) wmedianBottomTop(layers map[int]*graph.Layer) {
	medians := graph.NodeFloatMap{}
	for r := len(layers) - 1; r >= 0; r-- {
		for _, v := range layers[r].Nodes {
			medians[v] = medianOf(p.adjacentNodesPositions(v, v.Out, r+1))
		}
		p.sortLayer(layers[r].Nodes, medians)
	}
}

// The median of each vertex is the median of the positions of adjacent nodes in the previous (or following) layer.
func medianOf(adpos []int) float64 {
	// convert positions to float64 to simplify arithmetic ops
	fpos := make([]float64, len(adpos))
	for i, x := range adpos {
		fpos[i] = float64(x)
	}

	mid := len(fpos) / 2

	// Nodes with no adjacent vertices are given a median value of -1. This is used within the sort
	// function to indicate that these nodes should be left in their current positions.
	switch {
	case len(fpos) == 0: // no adjacent nodes
		return -1.0

	case len(fpos)%2 == 1: // odd number of adjacent nodes, get the median value
		return fpos[mid]

	case len(fpos) == 2: // get average of the two values
		return (fpos[0] + fpos[1]) / 2

	default:
		left := fpos[mid-1] - fpos[0]
		right := fpos[len(fpos)-1] - fpos[mid]
		return (fpos[mid-1]*right + fpos[mid]*left) / (left + right)
	}
}

// returns an ordered array of the present positions of the nodes
// adjacent to v in the given adjacent rank.
func (p *graphvizDotProcessor) adjacentNodesPositions(n *graph.Node, edges []*graph.Edge, adjLayer int) []int {
	res := []int{}
	for _, e := range edges {
		if e.SelfLoops() {
			continue
		}
		m := e.ConnectedNode(n)
		if m.Layer == adjLayer {
			res = append(res, p.getPos(m))
		}
	}
	sort.Ints(res)
	return res
}

func (p *graphvizDotProcessor) sortLayer(nodes []*graph.Node, medians graph.NodeFloatMap) {
	sort.Slice(nodes, func(i, j int) bool {
		a, aitr := p.fixedPositions.head(nodes[i])
		b, bitr := p.fixedPositions.head(nodes[j])

		if (aitr != 0 || bitr != 0) && a == b {
			return aitr < bitr
		}
		a_before_b := p.getPos(a) < p.getPos(b)
		b_before_a := p.getPos(b) < p.getPos(a)
		afixed := medians[a] == -1 && a_before_b
		bfixed := medians[b] == -1 && b_before_a

		flipIfNotInClosure := p.flipEqual && aitr == 0 && bitr == 0 && medians[a] == medians[b] && b_before_a

		return afixed || bfixed || medians[a] < medians[b] || flipIfNotInClosure
	})
	for i, n := range nodes {
		p.setPos(n, i)
	}
}

// transpose sweeps through layers in order and swaps pairs of adjacent nodes in the same layer;
// it counts the number of crossings between L, L-1 and L+1, if there's an improvement it keeps looping
// until no improvement is found.
func (p *graphvizDotProcessor) transpose(layers map[int]*graph.Layer) {
	improved := true
	for improved {
		improved = false
		for L := 0; L < len(layers); L++ {
			for i := 0; i < len(layers[L].Nodes)-2; i++ {
				v := layers[L].Nodes[i]
				w := layers[L].Nodes[i+1]

				// todo: the no-flip logic based on flat edges can be improved to consider
				// 	the closure as if it were one node. the non-closure node could be brought to the other end
				// 	of the closure

				// if w is head, skip
				if p.fixedPositions.mustBefore[v] == nil && p.fixedPositions.mustBefore[w] != nil {
					continue
				}
				// if v is tail, skip
				if p.fixedPositions.mustAfter[v] != nil && p.fixedPositions.mustAfter[w] == nil {
					continue
				}

				curX := crossingsAround(L, layers)
				p.swap(v, w)
				newX := crossingsAround(L, layers)

				switch {
				case newX < curX:
					// improved and keep new order
					improved = true
					layers[L].Nodes[i] = w
					layers[L].Nodes[i+1] = v

				default:
					// no improvement, restore order
					p.swap(v, w)
				}
			}
		}
	}
}

func crossings(layers map[int]*graph.Layer) int {
	crossings := 0
	for l := 1; l < len(layers); l++ {
		crossings += layers[l].CountCrossings()
	}
	return crossings
}

func crossingsAround(l int, layers map[int]*graph.Layer) int {
	if l == 0 {
		return layers[l+1].CountCrossings()
	}
	if l == len(layers)-1 {
		return layers[l].CountCrossings()
	}
	return layers[l].CountCrossings() + layers[l+1].CountCrossings()
}

func (p *graphvizDotProcessor) swap(v, w *graph.Node) {
	iv := p.getPos(v)
	iw := p.getPos(w)
	p.setPos(v, iw)
	p.setPos(w, iv)
}

func (p *graphvizDotProcessor) getPos(n *graph.Node) int {
	pos := p.positions[n]
	if pos != n.LayerPos {
		panic("gansner-north orderer: corrupted state: node in-layer position mismatch")
	}
	return pos
}

func (p *graphvizDotProcessor) setPos(n *graph.Node, pos int) {
	p.positions[n] = pos
	n.LayerPos = pos
}

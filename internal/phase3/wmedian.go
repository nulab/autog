package phase3

import (
	"sort"

	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

type initDirection uint8

const (
	initDirectionTop    initDirection = 0
	initDirectionBottom initDirection = 1
)

type wmedianProcessor struct {
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
func execWeightedMedian(g *graph.DGraph, params graph.Params) {
	if len(g.Layers) == 1 {
		// no crossings to reduce
		return
	}

	// insert virtual nodes so that edges with length >1 have length 1
	breakLongEdges(g)
	// set size to virtual nodes if needed
	if s := params.VirtualNodeFixedSize; s > 0.0 {
		for n := range g.VirtualNodes() {
			n.W = s
			n.H = s
		}
	}

	maxiter := int(params.WMedianMaxIter)
	fixedPositions := initFixedPositions(g.Edges)

	bestx_top, bestpos_top := wmedianRun(g, wmedianRunParams{maxiter, fixedPositions, initDirectionTop})
	bestx_btm, bestpos_btm := wmedianRun(g, wmedianRunParams{maxiter, fixedPositions, initDirectionBottom})

	var (
		bestx                  = 0
		bestp graph.NodeIntMap = nil
	)
	if bestx_top < bestx_btm {
		bestx, bestp = bestx_top, bestpos_top
	} else {
		bestx, bestp = bestx_btm, bestpos_btm
	}

	imonitor.Log("crossings", bestx)

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

type wmedianRunParams struct {
	maxiter        int
	fixedPositions fixedPositions
	dir            initDirection
}

type wmedianInitFn func(n *graph.Node, visited graph.NodeSet, indices map[int]int)

// node order is maintained in three different places:
//   - in g.Layers.Nodes, which is a slice
//   - in each node.LayerPos field
//   - in p.positions
//
// at each iteration, this algorithm will update the node positions in all three places
// a copy of the best p.positions is kept and at the end it is propagated to g.Layers and node.LayerPos
func wmedianRun(g *graph.DGraph, params wmedianRunParams) (int, graph.NodeIntMap) {
	p := &wmedianProcessor{
		positions:      graph.NodeIntMap{},
		fixedPositions: params.fixedPositions,
	}
	switch params.dir {
	case initDirectionTop:
		p.initPositions(g, g.Layers[0], p.initPositionsFromTop)
	case initDirectionBottom:
		p.initPositions(g, g.Layers[len(g.Layers)-1], p.initPositionsFromBottom)
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
	if bestx == 0 {
		return bestx, bestp
	}

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

func (p *wmedianProcessor) initPositions(g *graph.DGraph, layer *graph.Layer, fn wmedianInitFn) {
	// initialize positions
	visited := graph.NodeSet{}
	indices := map[int]int{}
	for _, n := range layer.Nodes {
		fn(n, visited, indices)
	}
	for _, n := range g.Nodes {
		fn(n, visited, indices)
	}
}

func (p *wmedianProcessor) initPositionsFromTop(n *graph.Node, visited graph.NodeSet, indices map[int]int) {
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

func (p *wmedianProcessor) initPositionsFromBottom(n *graph.Node, visited graph.NodeSet, indices map[int]int) {
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

func (p *wmedianProcessor) initPositionsFlatEdges(n *graph.Node, visited graph.NodeSet, indices map[int]int) {
	h, i := p.fixedPositions.head(n)
	if i > 0 {
		for h != nil && h != n {
			if visited[h] {
				goto next
			}
			visited[h] = true
			p.setPos(h, indices[h.Layer])
			indices[h.Layer]++
		next:
			h = p.fixedPositions.mustBefore[h]
		}
	}
}

// The weighted median routine assigns an order to each vertex in layer L(i) based on the current order
// of adjacent nodes in the next rank. Next is L(i)-1 in top-bottom sweep, or L(i)+1 in bottom-top sweep.
// Nodes with no adjacent nodes in the next layer are kept in place.
func (p *wmedianProcessor) wmedianTopBottom(layers []*graph.Layer) {
	medians := graph.NodeFloatMap{}
	for r := 1; r < len(layers); r++ {
		for _, v := range layers[r].Nodes {
			medians[v] = medianOf(p.adjacentNodesPositions(v, v.In, r-1))
		}
		p.sortLayer(layers[r].Nodes, medians)
	}
}

func (p *wmedianProcessor) wmedianBottomTop(layers []*graph.Layer) {
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
		if left == right {
			return (fpos[mid-1] + fpos[mid]) / 2
		}
		return (fpos[mid-1]*right + fpos[mid]*left) / (left + right)
	}
}

// returns an ordered array of the present positions of the nodes
// adjacent to v in the given adjacent rank.
func (p *wmedianProcessor) adjacentNodesPositions(n *graph.Node, edges []*graph.Edge, adjLayer int) []int {
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

func (p *wmedianProcessor) sortLayer(nodes []*graph.Node, medians graph.NodeFloatMap) {
	ep := len(nodes) // end pointer
	// back iteration slightly more efficient because it compares the iteration variable to zero
	for iter := len(nodes) - 1; iter >= 0; iter-- {
		lp := 0 // left pointer
		rp := 0 // right pointer
		for lp < ep {
			for lp < ep && medians[nodes[lp]] == -1 {
				lp++
			}
			if lp >= ep {
				break
			}
			canSwap := true
			for rp = lp + 1; rp < ep; rp++ {
				// todo: this probably isn't enough to handle flat edges
				if h, i := p.fixedPositions.head(nodes[rp]); i > 0 && h == nodes[lp] {
					canSwap = false
					break
				}
				if medians[nodes[rp]] >= 0 {
					break
				}
			}
			if rp >= ep {
				// no swap candidate
				break
			}
			if canSwap {
				ml := medians[nodes[lp]]
				mr := medians[nodes[rp]]
				if ml > mr || (ml == mr && p.flipEqual) {
					p.swap(nodes[lp], nodes[rp])
					nodes[lp], nodes[rp] = nodes[rp], nodes[lp]
				}
			}
			lp = rp
		}
		// todo: hasfixed in dot is true if there are flat edges in this layer
		if /*!hasfixed &&*/ !p.flipEqual {
			ep--
		}
	}
}

// transpose sweeps through layers in order and swaps pairs of adjacent nodes in the same layer;
// it counts the number of crossings between L, L-1 and L+1, if there's an improvement it keeps looping
// until no improvement is found.
func (p *wmedianProcessor) transpose(layers []*graph.Layer) {
	improved := true
	for improved {
		improved = false
		for _, layer := range layers {
			for i := 0; i < len(layer.Nodes)-2; i++ {
				v := layer.Nodes[i]
				w := layer.Nodes[i+1]

				if p.fixedPositions.mustBefore[v] == w {
					continue
				}

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

				curX := crossingsAround(layer.Index, layers)
				p.swap(v, w)
				newX := crossingsAround(layer.Index, layers)

				switch {
				case newX < curX:
					// improved and keep new order
					improved = true
					layer.Nodes[i] = w
					layer.Nodes[i+1] = v

				default:
					// no improvement, restore order
					p.swap(v, w)
				}
			}
		}
	}
}

func crossings(layers []*graph.Layer) int {
	crossings := 0
	for l := 1; l < len(layers); l++ {
		crossings += countCrossings(layers[l-1], layers[l])
	}
	return crossings
}

func crossingsAround(l int, layers []*graph.Layer) int {
	if l == 0 {
		return countCrossings(layers[l], layers[l+1])
	}
	if l == len(layers)-1 {
		return countCrossings(layers[l-1], layers[l])
	}
	return countCrossings(layers[l-1], layers[l]) + countCrossings(layers[l], layers[l+1])
}

func (p *wmedianProcessor) swap(v, w *graph.Node) {
	iv := p.getPos(v)
	iw := p.getPos(w)
	p.setPos(v, iw)
	p.setPos(w, iv)
}

func (p *wmedianProcessor) getPos(n *graph.Node) int {
	pos := p.positions[n]
	if pos != n.LayerPos {
		panic("gansner-north orderer: corrupted state: node in-layer position mismatch")
	}
	return pos
}

func (p *wmedianProcessor) setPos(n *graph.Node, pos int) {
	p.positions[n] = pos
	n.LayerPos = pos
}

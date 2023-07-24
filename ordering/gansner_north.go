package ordering

import (
	"math"
	"slices"
	"sort"
	"strconv"

	"github.com/vibridi/autog/graph"
)

const (
	maxiter = 24
	improv  = 0.02
)

type gansnerNorthProcessor struct {
	layers     map[int][]*graph.Node
	orders     graph.NodeMap
	minL, maxL int
	crossings  int
}

// // this implements a layered graph node ordering algorithm, based on:
// //   - "Emden R. Gansner, Eleftherios Koutsofios, Stephen C. North, Kiem-Phong Vo, A technique for
// //     drawing directed graphs. Software Engineering 19(3), pp. 214-230, 1993."
// //     https://www.researchgate.net/publication/3187542_A_Technique_for_Drawing_Directed_Graphs
func execGansnerNorth(g *graph.DGraph) {
	// insert virtual nodes so that edges with length >1 have length 1
	breakLongEdges(g)
	p := &gansnerNorthProcessor{
		layers: map[int][]*graph.Node{},
		orders: graph.NodeMap{},
		minL:   math.MaxInt,
		maxL:   math.MinInt,
	}

	for _, n := range g.Nodes {
		p.layers[n.Layer] = append(p.layers[n.Layer], n)
		p.minL = min(p.minL, n.Layer)
		p.maxL = max(p.maxL, n.Layer)
	}
	if len(p.layers) == 1 {
		return
	}

	visited := graph.NodeSet{}
	indices := map[int]int{}
	for _, n := range g.Nodes {
		p.initOrder(n, visited, indices)
	}
	for _, layer := range p.layers {
		slices.SortFunc(layer, func(a, b *graph.Node) int {
			if a.Layer != b.Layer {
				panic("same-layer nodes have different layers")
			}
			return p.orders[a] - p.orders[b]
		})
	}

	best := p.allCrossings()
	q := p.clone()

	for i := 0; i < maxiter; i++ {
		// Depending on the parity of the current iteration
		// number, the ranks are traversed from top to bottom or from bottom to top.
		if i%2 == 0 {
			p.wmedianTopBottom()
		} else {
			p.wmedianBottomTop()
		}
		p.transpose()
		crossings := p.allCrossings()
		if crossings < best {
			q = p.clone()
			best = crossings
		}
	}
	for _, n := range g.Nodes {
		n.LayerIdx = q.orders[n]
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
			}
			v++
			// set e's target to the virtual node
			e.To = virtualNode
			// add e to virtual node incoming edges
			virtualNode.In = append(virtualNode.In, e)
			// create new edge from virtual to e's former target
			f := &graph.Edge{From: virtualNode, To: to}
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
			// restart loop
			goto loop
		}
	}
}

// initially orders the nodes in each rank. This may be done by a depth-first or
// breadth-first search starting with vertices of minimum rank. Vertices are assigned positions in
// their ranks in left-to-right order as the search progresses. This strategy ensures that the initial
// ordering of a tree has no crossings. This is important because such crossings are obvious, easily avoided ‘‘mistakes.’’
func (p *gansnerNorthProcessor) initOrder(n *graph.Node, visited graph.NodeSet, indices map[int]int) {
	if visited[n] {
		return
	}
	visited[n] = true
	p.orders[n] = indices[n.Layer]
	indices[n.Layer]++
	for _, e := range n.Out {
		p.initOrder(e.To, visited, indices)
	}
}

// At each rank a vertex is assigned a median based on the adjacent vertices on the previous
// rank. Then, the vertices in the rank are sorted by their medians. An important consideration is
// what to do with vertices that have no adjacent vertices on the previous rank. In our
// implementation such vertices are left fixed in their current positions with non-fixed vertices sorted
// into the remaining positions.
func (p *gansnerNorthProcessor) wmedianTopBottom() {
	medians := map[*graph.Node]float64{}
	for r := p.minL + 1; r < p.maxL; r++ {
		for _, v := range p.layers[r] {
			medians[v] = p.medianOf(p.adjacentNodesPositions(v, v.In, r-1))
		}
		p.sortLayer(p.layers[r], medians)
	}
}

func (p *gansnerNorthProcessor) wmedianBottomTop() {
	medians := map[*graph.Node]float64{}
	for r := p.maxL - 1; r >= p.minL; r-- {
		for _, v := range p.layers[r] {
			medians[v] = p.medianOf(p.adjacentNodesPositions(v, v.Out, r+1))
		}
		p.sortLayer(p.layers[r], medians)
	}
}

// The median value of a vertex is defined as the median position of the adjacent vertices if that
// is uniquely defined. Otherwise, it is interpolated between the two median positions using a
// measure of tightness. Generally, the weighted median is biased toward the side where vertices are
// more closely packed.
func (p *gansnerNorthProcessor) medianOf(adpos []int) float64 {
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
func (p *gansnerNorthProcessor) adjacentNodesPositions(n *graph.Node, edges []*graph.Edge, adjLayer int) []int {
	res := []int{}
	for _, e := range edges {
		if e.SelfLoops() {
			continue
		}
		m := e.ConnectedNode(n)
		if m.Layer == adjLayer {
			res = append(res, p.orders[m])
		}
	}
	sort.Ints(res)
	return res
}

func (p *gansnerNorthProcessor) sortLayer(layer []*graph.Node, medians map[*graph.Node]float64) {
	slices.SortFunc(layer, func(a, b *graph.Node) int {
		afixed := medians[a] == -1 && p.orders[a] < p.orders[b]
		bfixed := medians[b] == -1 && p.orders[b] < p.orders[a]
		if afixed || bfixed || medians[a] < medians[b] {
			return -1
		}
		return 1
	})
	for i, n := range layer {
		p.orders[n] = i
	}
}

// 3-15: This is the main loop that iterates as long as the number of edge crossings can be reduced by
// transpositions. As in the loop in the ordering function, an adaptive strategy could be applied
// here to terminate the loop once the improvement is a sufficiently small fraction of the number of
// crossings.
// 7-12: Each adjacent pair of vertices is examined. Their order is switched if this reduces the number of
// crossings. The function crossing(v,w) simply counts the number of edge crossings if v
// appears to the left of w in their rank
func (p *gansnerNorthProcessor) transpose() {
	// todo: adaptive strategy to keep iterating in case of sufficiently large improvement
	for improved, itr := true, 0; improved && itr < 20; itr++ {
		improved = false
		for L := p.minL; L <= p.maxL; L++ {
			for i := 0; i < len(p.layers[L])-2; i++ {
				v := p.layers[L][i]
				w := p.layers[L][i+1]
				curX := p.layerCrossings(L, v, w)
				newX := p.layerCrossings(L, w, v)
				if curX > newX {
					improved = true
					p.swapOrder(v, w)
					p.layers[L][i] = w
					p.layers[L][i+1] = v
				}
			}
		}
	}
}

func (p *gansnerNorthProcessor) allCrossings() int {
	crossings := 0
	for l := 1; l <= p.maxL; l++ {
		crossings += p.layerCrossings(l, nil, nil)
	}
	p.crossings = crossings
	return crossings
}

func (p *gansnerNorthProcessor) layerCrossings(l int, v, w *graph.Node) int {
	switch {
	case l == p.minL:
		return p.crossingsOf(l+1, v, w)
	case l == p.maxL:
		return p.crossingsOf(l, v, w)
	default:
		return p.crossingsOf(l, v, w) + p.crossingsOf(l+1, v, w)
	}
}

// determines the number of crossings between L and L-1 that involve v and w. crossings that aren't affected by v and w's relative
// order aren't counted.
// given that long edges have been broken by inserting virtual nodes, and that all edges
// connect nodes only one layer apart, crossings are determined purely by node order.
func (p *gansnerNorthProcessor) crossingsOf(l int, v, w *graph.Node) int {
	nodes := p.layers[l]
	// swap order
	if v != nil && w != nil && p.orders[v] > p.orders[w] {
		nodes = []*graph.Node{v, w}
		p.swapOrder(v, w)
		defer p.swapOrder(v, w) // undo the swap
	}
	crossings := 0
	visited := map[uint64]bool{}
	if l > p.minL {
		// examine crossings with upper layer: l-1
		for _, n := range nodes {
			for _, e := range n.In {
				if e.SelfLoops() {
					continue
				}
				upperLayer := p.layers[l-1]
				for i := p.orders[n]; i < len(upperLayer); i++ {
					for _, f := range upperLayer[i].Out {
						if f == e || visited[p.bitmask(e, f)] {
							continue
						}
						visited[p.bitmask(e, f)] = true
						if p.edgesCross(e, f) {
							crossings++
						}
					}
				}
			}
		}
	}
	return crossings
}

func (p *gansnerNorthProcessor) swapOrder(v, w *graph.Node) {
	iv := p.orders[v]
	iw := p.orders[w]
	p.orders[v] = iw
	p.orders[w] = iv
}

func (p *gansnerNorthProcessor) edgesCross(e, f *graph.Edge) bool {
	right2Left := p.orders[f.From] > p.orders[e.From] && p.orders[f.To] < p.orders[e.To]
	left2Right := p.orders[f.From] < p.orders[e.From] && p.orders[f.To] > p.orders[e.To]
	return right2Left || left2Right
}

func (p *gansnerNorthProcessor) bitmask(e, f *graph.Edge) uint64 {
	x := uint64(0)
	x |= 1 << p.nodeNum(e.From)
	x |= 1 << p.nodeNum(e.To)
	x |= 1 << p.nodeNum(f.From)
	x |= 1 << p.nodeNum(f.To)
	return x
}

func (p *gansnerNorthProcessor) nodeNum(n *graph.Node) uint64 {
	x := p.orders[n]
	for i := n.Layer - 1; i >= p.minL; i-- {
		x += len(p.layers[i])
	}
	return uint64(x)
}

func (p *gansnerNorthProcessor) clone() *gansnerNorthProcessor {
	layers := map[int][]*graph.Node{}
	for k, v := range p.layers {
		layers[k] = append([]*graph.Node{}, v...)
	}

	orders := graph.NodeMap{}
	for k, v := range p.orders {
		orders[k] = v
	}
	return &gansnerNorthProcessor{
		layers: layers,
		orders: orders,
		minL:   p.minL,
		maxL:   p.maxL,
	}
}

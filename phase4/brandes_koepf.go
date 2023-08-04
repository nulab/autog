package phase4

import (
	"fmt"
	"maps"
	"math"
	"sort"

	"github.com/nulab/autog/graph"
)

// todo: make sure this accounts for node sizes and ports.
// Rueegg-Schulze https://rtsys.informatik.uni-kiel.de/~biblio/downloads/papers/gd15.pdf
// If ports aren't relevant to a particular implementation, node size still is, so the port can be set by default
// at the middle point of the node side.

type direction uint8

const (
	top direction = iota
	bottom
	left
	right
)

type layout struct {
	v, h      direction
	blockroot graph.NodeMap
	alignment graph.NodeMap
}

type pair struct {
	node *graph.Node
	edge *graph.Edge
}

type xcoordinates graph.NodeFloatMap

func (xc xcoordinates) Size() (w, minx, maxx float64) {
	minx = math.Inf(+1)
	maxx = math.Inf(-1)
	for n, x := range xc {
		minx = min(minx, x)
		maxx = max(maxx, x+n.W)
	}
	w = maxx - minx
	return
}

type brandesKoepfPositioner struct {
	markedEdges graph.EdgeSet
	neighbors   map[*graph.Node]map[direction][]pair
	layerFor    func(*graph.Node) *graph.Layer
	nodeMargin  float64
	nodeSpacing float64
}

// this implements an O(n) time x-coordinate assignment algorithm, based on:
//   - "Ulrik Brandes and Boris Köpf, Fast and Simple Horizontal Coordinate Assignment"
//     https://link.springer.com/content/pdf/10.1007/3-540-45848-4_3.pdf
//   - ELK Java code at https://github.com/eclipse/elk/tree/master/plugins/org.eclipse.elk.alg.layered/src/org/eclipse/elk/alg/layered/p4nodes/bk
func execBrandesKoepf(g *graph.DGraph, params graph.Params) {
	p := &brandesKoepfPositioner{
		markedEdges: graph.EdgeSet{},
		neighbors:   neighbors(g),
		nodeMargin:  params.NodeMargin,
		nodeSpacing: params.NodeSpacing,

		layerFor: func(n *graph.Node) *graph.Layer {
			return g.Layers[n.Layer]
		},
	}
	p.markConflicts(g)

	layouts := [4]layout{
		{v: bottom, h: right},
		{v: bottom, h: left},
		{v: top, h: right},
		{v: top, h: left},
	}
	xcoords := [4]xcoordinates{}

	for i, a := range layouts {
		// initialize per-layout blocks and alignment maps
		a.blockroot = make(graph.NodeMap, len(g.Nodes))
		a.alignment = make(graph.NodeMap, len(g.Nodes))
		for _, n := range g.Nodes {
			a.blockroot[n] = n
			a.alignment[n] = n
		}
		// main phases
		p.verticalAlign(g, a)
		xcoords[i] = p.horizontalCompaction(g, a)
	}

	finalLayout := balanceLayouts(xcoords, g.Nodes)

	if !verifyLayout(finalLayout, g.Layers, params.NodeMargin) {
		changed := false
		smallest, _, _ := finalLayout.Size()

		for _, xc := range xcoords {
			if verifyLayout(xc, g.Layers, params.NodeMargin) {
				if w, _, _ := xc.Size(); w < smallest {
					smallest = w
					finalLayout = xc
					changed = true
				}
			}
		}
		if !changed {
			fmt.Println("NO VIABLE LAYOUT")
			// finalLayout = xcoords[0]
		}
	}

	for _, l := range g.Layers {
		for _, n := range l.Nodes {
			n.X = finalLayout[n] // + (params.NodeSpacing+params.NodeMargin*2)*float64(i)
			l.H = max(l.H, n.H)
		}
	}
}

// marks edges that cross inner edges, i.e. type 1 and type 2 conflicts as defined in B&K
func (p *brandesKoepfPositioner) markConflicts(g *graph.DGraph) {
	if len(g.Layers) < 4 {
		return
	}
	// sweep layers from top to bottom except the first and the last
	for i := 1; i < len(g.Layers)-1; i++ {
		k0 := 0
		for l1, v := range g.Layers[i+1].Nodes {
			ksrc := incidentToInner(v)
			if g.Layers[i+1].Tail() == v || ksrc >= 0 {
				// set k1 to the index of the second-to-last node or the previous, or
				// if v belongs to an inner edge, to the leftmost upper neighbor of v
				k1 := g.Layers[i].Len() - 1
				if ksrc >= 0 {
					k1 = p.neighbors[v][bottom][0].node.LayerPos
				}
				// range over same layer nodes until v included
				for l2, w := range g.Layers[i+1].Nodes {
					if l2 > l1 {
						break
					}
					for _, e := range w.In {
						if e.SelfLoops() || e.IsFlat() {
							continue
						}
						if e.From.LayerPos < k0 || e.From.LayerPos > k1 {
							p.markedEdges[e] = true
						}
					}
				}
				k0 = k1
			}
		}
	}
}

func neighbors(g *graph.DGraph) map[*graph.Node]map[direction][]pair {
	neighbors := map[*graph.Node]map[direction][]pair{}
	for _, n := range g.Nodes {
		neighbors[n] = make(map[direction][]pair, 2)
		if n.Layer > 0 {
			var ps []pair
			for _, e := range n.In {
				if !e.SelfLoops() && !e.IsFlat() {
					ps = append(ps, pair{e.From, e})
				}
			}
			// when sweeping layers downward, we want to examine upper neighbors
			neighbors[n][bottom] = ps
		}

		if n.Layer < len(g.Layers)-1 {
			var ps []pair
			for _, e := range n.Out {
				if !e.SelfLoops() && !e.IsFlat() {
					ps = append(ps, pair{e.To, e})
				}
			}
			// when sweeping layers upward, we want to examine lower neighbors
			neighbors[n][top] = ps
		}
	}
	return neighbors
}

func (p *brandesKoepfPositioner) verticalAlign(g *graph.DGraph, layout layout) {
	iter := layersIterator(g, layout.v)
	for layer := iter(); layer != nil; layer = iter() {
		// r is the index of the nearest neighbor to which vk can be aligned
		// by updating r with the most recently aligned neighbor (at the end of the loop)
		// it's guaranteed that only one alignment is possible
		r := outermostPos(layout.h)
		iter := nodesIterator(layer.Nodes, layout.h)
		for vk := iter(); vk != nil; vk = iter() {
			vkneighbors := p.neighbors[vk][layout.v]
			if d := len(vkneighbors); d > 0 {
				for _, m := range medianNeighborIndices(d, layout.h) {
					if layout.alignment[vk] == vk /* not aligned */ {
						u, uv := vkneighbors[m].node, vkneighbors[m].edge
						if !p.markedEdges[uv] && withinOutermostPos(r, u.LayerPos, layout.h) {
							// align and blockroot maintain a circular reference:
							// in top-bottom direction, a node u aligns with a lower one vk
							// and vk aligns with the root of its block
							layout.alignment[u] = vk
							layout.blockroot[vk] = layout.blockroot[u]
							layout.alignment[vk] = layout.blockroot[vk]
							r = u.LayerPos
						}
					}
				}
			}
		}
	}
}

type classes struct {
	sinks  graph.NodeMap // sink blocks
	xshift graph.NodeFloatMap
	xcoord graph.NodeFloatMap
	xcinit graph.NodeSet
}

func (p *brandesKoepfPositioner) horizontalCompaction(g *graph.DGraph, layout layout) xcoordinates {
	c := &classes{
		sinks:  graph.NodeMap{},
		xshift: graph.NodeFloatMap{},
		xcoord: graph.NodeFloatMap{},
		xcinit: graph.NodeSet{},
	}
	for _, n := range g.Nodes {
		c.sinks[n] = n
		c.xshift[n] = outermostX(layout.h)
	}

	iter := layersIterator(g, layout.v)
	for layer := iter(); layer != nil; layer = iter() {
		iter := nodesIterator(layer.Nodes, layout.h)
		for n := iter(); n != nil; n = iter() {
			if layout.blockroot[n] == n {
				p.placeBlock(n, c, layout)
			}
		}
	}

	for _, n := range g.Nodes {
		c.xcoord[n] = c.xcoord[layout.blockroot[n]]
		if shift := c.xshift[c.sinks[layout.blockroot[n]]]; withinOutermostX(shift, layout.h) {
			c.xcoord[n] = c.xcoord[n] + shift
		}
	}
	return xcoordinates(c.xcoord)
}

func (p *brandesKoepfPositioner) placeBlock(v *graph.Node, c *classes, layout layout) {
	if c.xcinit[v] {
		// already placed
		return
	}
	c.xcinit[v] = true
	c.xcoord[v] = 0
	w := v
	for {
		wlayer := p.layerFor(w)
		leftNotLast := layout.h == left && w.LayerPos > 0
		rightNotLast := layout.h == right && w.LayerPos < len(wlayer.Nodes)-1

		if leftNotLast || rightNotLast {
			u := previousNodeInLayer(w, wlayer.Nodes, layout.h)
			uroot := layout.blockroot[u]
			p.placeBlock(uroot, c, layout)
			if c.sinks[v] == v {
				c.sinks[v] = c.sinks[uroot]
			}
			if c.sinks[v] != c.sinks[uroot] {
				switch layout.h {
				case left:
					c.xshift[c.sinks[uroot]] = min(c.xshift[c.sinks[uroot]], c.xcoord[v]-c.xcoord[uroot]-p.space(u))
				case right:
					c.xshift[c.sinks[uroot]] = max(c.xshift[c.sinks[uroot]], c.xcoord[v]-c.xcoord[uroot]+p.space(u))
				}
			} else {
				switch layout.h {
				case left:
					c.xcoord[v] = max(c.xcoord[v], c.xcoord[uroot]+p.space(u))
				case right:
					c.xcoord[v] = min(c.xcoord[v], c.xcoord[uroot]-p.space(u))
				}

			}
		}
		// the align map contains the next node in the block
		w = layout.alignment[w]
		if w == v {
			// back at root
			break
		}
	}
}

// type block = []*graph.Edge
//
// var inn = graph.NodeFloatMap{}
// var blox = [][]*graph.Edge{}
// var blockSize = map[*block]float64{}
//
// func innerShift(nodes []*graph.Node) {
// 	for _, n := range nodes {
// 		inn[n] = 0
// 		for _, b := range blox {
// 			left, right := 0.0, 0.0
// 			for _, e := range b {
// 				p, q := e.From, e.To
// 				s := inn[π(p)] + xp(p) - xp(q)
// 				inn[π(q)] = s
// 				left = min(left, s)
// 				right = max(right, s+width(π(q)))
// 			}
// 			for _, e := range blox {
// 				n := (*graph.Node)(unsafe.Pointer(e[0]))
// 				inn[n] -= left
// 			}
// 			blockSize[&b] = right - left
// 		}
// 	}
// }
//
// type port = *graph.Node // todo
//
// // maps port to node
// func π(port) *graph.Node {
// 	return nil
// }
//
// func xp(port) float64 {
// 	return 0.0
// }
//
// func width(port) float64 {
// 	return 0
// }

// returns a non-negative number if n is the target node of an inner edge, i.e. an edge connecting two virtual nodes
// on adjacent layers, where the number is the position of the edge source in the upper layer;
// it returns -1 if the node isn't involved in an inner edge.
func incidentToInner(n *graph.Node) int {
	if !n.IsVirtual {
		return -1
	}
	for _, e := range n.In {
		if e.From.IsVirtual && e.From.Layer == n.Layer-1 {
			return e.From.LayerPos
		}
	}
	return -1
}

func layersIterator(g *graph.DGraph, dir direction) func() *graph.Layer {
	ks := maps.Keys(g.Layers)
	switch dir {
	case bottom:
		sort.Ints(ks)
	case top:
		sort.Sort(sort.Reverse(sort.IntSlice(ks)))
	default:
		panic("BK positioner: invalid layer iteration direction")
	}
	i := 0
	return func() *graph.Layer {
		if i >= len(ks) {
			return nil
		}
		layer := g.Layers[ks[i]]
		i++
		return layer
	}
}

func nodesIterator(nodes []*graph.Node, dir direction) func() *graph.Node {
	var i int
	switch dir {
	case right:
		i = 0
	case left:
		i = len(nodes) - 1
	default:
		panic("BK positioner: invalid node iteration direction")
	}
	return func() *graph.Node {
		if (dir == right && i >= len(nodes)) || (dir == left && i < 0) {
			return nil
		}
		node := nodes[i]
		switch dir {
		case right:
			i++
		case left:
			i--
		}
		return node
	}
}

func medianNeighborIndices(d int, dir direction) []int {
	// remember that indices in the paper start at 1 but here start at 0
	m1 := int(math.Floor((float64(d)+1.0)/2.0)) - 1
	m2 := int(math.Ceil((float64(d)+1.0)/2.0)) - 1
	switch dir {
	case right:
		return []int{m1, m2}
	case left:
		return []int{m2, m1}
	default:
		panic("BK positioner: invalid horizontal direction")
	}
}

func outermostPos(dir direction) int {
	switch dir {
	case right:
		return -1
	case left:
		return math.MaxInt
	default:
		panic("BK positioner: invalid horizontal direction")
	}
}

func outermostX(dir direction) float64 {
	switch dir {
	case right:
		return math.Inf(+1)
	case left:
		return math.Inf(-1)
	default:
		panic("BK positioner: invalid vertical direction")
	}
}

func withinOutermostPos(r, pos int, dir direction) bool {
	switch dir {
	case right:
		return r < pos
	case left:
		return r > pos
	default:
		panic("BK positioner: invalid horizontal direction")
	}
}

func withinOutermostX(shift float64, dir direction) bool {
	switch dir {
	case right:
		return shift < math.Inf(+1)
	case left:
		return shift > math.Inf(-1)
	default:
		panic("BK positioner: invalid horizontal direction")
	}
}

func previousNodeInLayer(n *graph.Node, nodes []*graph.Node, dir direction) *graph.Node {
	switch dir {
	case right:
		return nodes[n.LayerPos+1]
	case left:
		return nodes[n.LayerPos-1]
	default:
		panic("BK positioner: invalid horizontal direction")
	}
}

func (p *brandesKoepfPositioner) space(n *graph.Node) float64 {
	return n.W + p.nodeSpacing + p.nodeMargin*2
}

func balanceLayouts(layoutXCoords [4]xcoordinates, nodes []*graph.Node) xcoordinates {
	minx := [4]float64{}
	maxx := [4]float64{}
	width := [4]float64{}

	leastWidth := 0

	for i, xc := range layoutXCoords {
		width[i], minx[i], maxx[i] = xc.Size()

		if width[leastWidth] > width[i] {
			leastWidth = i
		}
	}

	shift := [4]float64{}
	for i := range layoutXCoords {
		if i == 1 || i == 3 /* left */ {
			shift[i] = minx[leastWidth] - minx[i]
		} else {
			shift[i] = maxx[leastWidth] - maxx[i]
		}
	}

	medianx := xcoordinates{}
	for _, n := range nodes {
		xs := make([]float64, 4)
		for i := range layoutXCoords {
			xs[i] = layoutXCoords[i][n] + shift[i]
		}
		sort.Float64s(xs)
		medianx[n] = (xs[1] + xs[2]) / 2.0
	}
	return medianx
}

// todo: account for node margin/spacing
func verifyLayout(layout xcoordinates, layers map[int]*graph.Layer, nodeMargin float64) bool {
	for _, layer := range layers {
		pos := math.Inf(-1)
		for _, n := range layer.Nodes {
			left := layout[n] - nodeMargin
			right := layout[n] + n.W + nodeMargin

			if left > pos && right > pos {
				pos = right
			} else {
				return false
			}
		}
	}
	return true
}

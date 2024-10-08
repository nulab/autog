package phase4

import (
	"iter"
	"math"
	"slices"
	"sort"

	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
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

type (
	layout struct {
		v, h      direction
		blockroot graph.NodeMap
		alignment graph.NodeMap
	}
	pair struct {
		node *graph.Node
		edge *graph.Edge
	}
	xcoordinates graph.NodeFloatMap
	neighbors    map[*graph.Node]map[direction][]pair
)

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
	neighbors   neighbors
	layerFor    func(*graph.Node) *graph.Layer
	nodeSpacing float64
}

// this implements an O(n) time x-coordinate assignment algorithm, based on:
//   - "Ulrik Brandes and Boris Köpf, Fast and Simple Horizontal Coordinate Assignment"
//     https://link.springer.com/content/pdf/10.1007/3-540-45848-4_3.pdf
//   - "Brandes, Walter and Zink, Erratum: Fast and Simple Horizontal Coordinate Assignment"
//     https://arxiv.org/pdf/2008.01252
//   - ELK Java code at https://github.com/eclipse/elk/tree/master/plugins/org.eclipse.elk.alg.layered/src/org/eclipse/elk/alg/layered/p4nodes/bk
//
// note that ELK implements the Rüegg-Schulze extension to the original algorithm.
func execBrandesKoepf(g *graph.DGraph, params graph.Params) {
	neighbors := initNeighbors(g)
	markedEdges := markConflicts(g, neighbors)

	p := &brandesKoepfPositioner{
		markedEdges: markedEdges,
		neighbors:   neighbors,
		nodeSpacing: params.NodeSpacing,

		layerFor: func(n *graph.Node) *graph.Layer {
			return g.Layers[n.Layer]
		},
	}

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

	var finalLayout xcoordinates
	// override default choice if user specified a preference
	// otherwise balance the four layouts together
	forceLayout := params.BrandesKoepfLayout >= 0 && params.BrandesKoepfLayout < 4
	if forceLayout {
		finalLayout = xcoords[params.BrandesKoepfLayout]
	} else {
		// TODO: with the final step that ensures no overlaps, this verification step might be not needed any more
		finalLayout = balanceLayouts(xcoords, g.Nodes)
		if !verifyLayout(finalLayout, g.Layers, params.NodeSpacing) {
			changed := false
			smallest, _, _ := finalLayout.Size()

			for _, xc := range xcoords {
				if verifyLayout(xc, g.Layers, params.NodeSpacing) {
					if w, _, _ := xc.Size(); w < smallest {
						smallest = w
						finalLayout = xc
						changed = true
					}
				}
			}
			if !changed {
				// keep balanced layout
				imonitor.Log("layout verification", "no viable layout, keep balanced")
			}
		}
	}

	lmargin := 0.0
	for _, l := range g.Layers {
		for _, n := range l.Nodes {
			n.X = finalLayout[n]
			lmargin = min(lmargin, n.X)
			l.H = max(l.H, n.H)
		}
	}
	// normalize negative xs
	if lmargin < 0 {
		lmargin = math.Abs(lmargin)
		for _, n := range g.Nodes {
			n.X += lmargin
		}
	}

	// B&K could produce a positioning with overlaps after averaging
	// this is a final adjustment step to mitigate the issue
	for _, l := range g.Layers {
		for j := 1; j < l.Len(); j++ {
			v := l.Nodes[j-1]
			w := l.Nodes[j]

			overlaps := w.X > v.X && w.X < v.X+v.W
			if overlaps {
				shift := v.X + v.W + p.nodeSpacing - w.X
				w.X += shift
			}
		}
	}
}

// marks edges that cross inner edges, i.e. type 1 and type 2 conflicts as defined in B&K
func markConflicts(g *graph.DGraph, neighbors neighbors) graph.EdgeSet {
	markedEdges := graph.EdgeSet{}
	if len(g.Layers) < 4 {
		return markedEdges
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
					k1 = neighbors[v][bottom][0].node.LayerPos
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
							markedEdges[e] = true
						}
					}
				}
				k0 = k1
			}
		}
	}
	return markedEdges
}

func initNeighbors(g *graph.DGraph) neighbors {
	// use the unnamed type here so that subsequent make/append calls are more easily understood
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
			neighbors[n][bottom] = slices.Clip(ps)
		}

		if n.Layer < len(g.Layers)-1 {
			var ps []pair
			for _, e := range n.Out {
				if !e.SelfLoops() && !e.IsFlat() {
					ps = append(ps, pair{e.To, e})
				}
			}
			// when sweeping layers upward, we want to examine lower neighbors
			neighbors[n][top] = slices.Clip(ps)
		}
	}
	return neighbors
}

func (p *brandesKoepfPositioner) verticalAlign(g *graph.DGraph, layout layout) {
	for layer := range iterLayers(g.Layers, layout.v) {
		// r is the index of the nearest neighbor to which vk can be aligned
		// by updating r with the most recently aligned neighbor (at the end of the loop)
		// it's guaranteed that only one alignment is possible
		r := outermostPos(layout.h)
		for vk := range iterNodes(layer.Nodes, layout.h) {
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

	for layer := range iterLayers(g.Layers, layout.v) {
		for n := range iterNodes(layer.Nodes, layout.h) {
			if layout.blockroot[n] == n {
				p.placeBlock(n, c, layout)
			}
		}
	}

	for layer := range iterLayers(g.Layers, layout.v) {
		n := firstNodeInLayer(layer, layout.h)
		if c.sinks[n] != n {
			continue
		}
		if shift := c.xshift[c.sinks[n]]; shift == outermostX(layout.h) {
			c.xshift[c.sinks[n]] = 0
		}

		k := 0
		j := layer.Index
		for j < len(g.Layers) && k < g.Layers[j].Len() {

			vjk := g.Layers[j].Nodes[k]
			v := vjk

			if c.sinks[v] != c.sinks[vjk] {
				break
			}

			for layout.alignment[v] != layout.blockroot[v] {
				v = layout.alignment[v]
				lv := g.Layers[v.Layer]
				if v != firstNodeInLayer(lv, layout.h) {
					u := prevNodeInLayer(v, lv.Nodes, layout.h)
					switch layout.h {
					case left:
						s := c.xshift[c.sinks[v]] + c.xcoord[v] + (c.xcoord[u] + u.W + p.nodeSpacing)
						c.xshift[c.sinks[u]] = max(c.xshift[c.sinks[u]], s)
					case right:
						s := c.xshift[c.sinks[v]] + c.xcoord[v] - (c.xcoord[u] + p.nodeSpacing)
						c.xshift[c.sinks[u]] = min(c.xshift[c.sinks[u]], s)
					}
				}
				j++
			}
			k = v.LayerPos + 1
		}
	}

	for _, n := range g.Nodes {
		if shift := c.xshift[c.sinks[n]]; withinOutermostX(shift, layout.h) {
			c.xcoord[n] += c.xshift[c.sinks[n]]
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
		isLast := w == lastNodeInLayer(wlayer, layout.h)

		if !isLast {
			u := nextNodeInLayer(w, wlayer.Nodes, layout.h)
			uroot := layout.blockroot[u]
			p.placeBlock(uroot, c, layout)
			if c.sinks[v] == v {
				c.sinks[v] = c.sinks[uroot]
			}
			if c.sinks[v] == c.sinks[uroot] {
				switch layout.h {
				case left:
					s := c.xcoord[uroot] + p.space(u)
					c.xcoord[v] = max(c.xcoord[v], s)
				case right:
					s := c.xcoord[uroot] - p.space(v)
					c.xcoord[v] = min(c.xcoord[v], s)
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

	for layout.alignment[w] != v {
		w = layout.alignment[w]
		c.xcoord[w] = c.xcoord[v]
		c.sinks[w] = c.sinks[v]
	}
}

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

func iterLayers(layers []*graph.Layer, dir direction) iter.Seq[*graph.Layer] {
	switch dir {
	case bottom:
		return slices.Values(layers)
	case top:
		return func(yield func(*graph.Layer) bool) {
			for _, n := range slices.Backward(layers) {
				if !yield(n) {
					return
				}
			}
		}
	default:
		panic("autog: B&K: invalid layer iteration direction")
	}
}

func iterNodes(nodes []*graph.Node, dir direction) iter.Seq[*graph.Node] {
	switch dir {
	case right:
		return slices.Values(nodes)
	case left:
		return func(yield func(*graph.Node) bool) {
			for _, n := range slices.Backward(nodes) {
				if !yield(n) {
					return
				}
			}
		}
	default:
		panic("autog: B&K: invalid node iteration direction")
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

func firstNodeInLayer(l *graph.Layer, dir direction) *graph.Node {
	switch dir {
	case right:
		return l.Head()
	case left:
		return l.Tail()
	default:
		panic("BK positioner: invalid horizontal direction")
	}
}

func nextNodeInLayer(n *graph.Node, nodes []*graph.Node, dir direction) *graph.Node {
	switch dir {
	case right:
		return nodes[n.LayerPos+1]
	case left:
		return nodes[n.LayerPos-1]
	default:
		panic("BK positioner: invalid horizontal direction")
	}
}

func prevNodeInLayer(n *graph.Node, nodes []*graph.Node, dir direction) *graph.Node {
	switch dir {
	case right:
		return nodes[n.LayerPos-1]
	case left:
		return nodes[n.LayerPos+1]
	default:
		panic("BK positioner: invalid horizontal direction")
	}
}

func lastNodeInLayer(l *graph.Layer, dir direction) *graph.Node {
	switch dir {
	case right:
		return l.Tail()
	case left:
		return l.Head()
	default:
		panic("BK positioner: invalid horizontal direction")
	}
}

func (p *brandesKoepfPositioner) space(n *graph.Node) float64 {
	return n.W + p.nodeSpacing
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

// todo:
// it could be worth it to allow for some slack here by considering valid layouts where the nodes
// don't overlap with a fraction of the spacing between them, instead of mandating full spacing
// however on failure ELK returns the first layout, we still return the average.
func verifyLayout(layout xcoordinates, layers []*graph.Layer, nodeSpacing float64) bool {
	for _, layer := range layers {
		pos := math.Inf(-1)
		for _, n := range layer.Nodes {
			left := layout[n]
			right := layout[n] + n.W + nodeSpacing

			if left > pos && right > pos {
				pos = right
			} else {
				return false
			}
		}
	}
	return true
}

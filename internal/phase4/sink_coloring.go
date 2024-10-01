package phase4

import (
	"slices"

	"github.com/nulab/autog/internal/graph"
)

// This implements an heuristic for x-coordinate assignment similar to the one proposed by Brandes & Köpf.
// Vertices are partitioned in blocks (colored) based on vertical alignment and giving priority to long edges.
// The partitioning proceeds greedily sweeping the layers from bottom to top.
// Initial x-coordinates are assigned to all nodes by packing them as left as possible in their layer, while respecting
// their relative position. Finally, all nodes are assigned the maximum x in their block.
// If there's any overlap, nodes are shifted right recursively.
func execSinkColoring(g *graph.DGraph, params graph.Params) {
	colors := graph.NodeMap{}
	roots := graph.NodeMap{}
	for _, n := range g.Nodes {
		colors[n] = n
		roots[n] = n
	}

	// init edge conflicts
	// unlike B&K here we consider edge priorities
	edgePriority := map[int][]*graph.Edge{}

	// paint nodes, and remember the maximum same-color block width, O(n)
	blockwidth := graph.NodeFloatMap{}

	for _, layer := range slices.Backward(g.Layers) {
		for _, n := range layer.Nodes {
			_, w := setColor(n, colors, roots, edgePriority)
			blockwidth[roots[n]] = max(blockwidth[roots[n]], w)
		}
	}

	// init coordinates by packing nodes to the left, O(n)
	xcoord := graph.NodeFloatMap{}
	for _, layer := range g.Layers {
		x := 0.0
		for _, n := range layer.Nodes {
			xcoord[n] = x
			x += blockwidth[roots[n]] + params.NodeSpacing
		}
	}

	// compute maximum x coord for each block, O(n)
	blockmax := graph.NodeFloatMap{}
	for n, x := range xcoord {
		blockmax[roots[n]] = max(blockmax[roots[n]], x)
	}

	lmax := 0
	for _, l := range g.Layers {
		lmax = max(lmax, l.Len())
	}

	spacing := params.NodeSpacing
	placeBlock(g, lmax, spacing, blockmax, blockwidth, xcoord, roots)

	for _, l := range g.Layers {
		for _, n := range l.Nodes {
			n.X = xcoord[n]
			l.H = max(l.H, n.H)
		}
	}
}

func setColor(n *graph.Node, colors graph.NodeMap, roots graph.NodeMap, priority map[int][]*graph.Edge) (*graph.Node, float64) {
	if colors[n] != n || len(n.In) == 0 {
		// already painted or source node
		return n, n.W
	}

	// candidate edge
	var e *graph.Edge

	// prioritize edges connecting to virtual nodes
	for _, f := range n.In {
		if f.ConnectedNode(n).IsVirtual {
			e = f
		}
	}

	// make sure the edge is viable, i.e. that connects to an upper neighbor with no conflicts
	i := 0
	for e == nil || e.SelfLoops() || e.IsFlat() {
		if i >= len(n.In) {
			// no viable edges, n is a block root (possibly with cardinality 1)
			return n, n.W
		}
		e = n.In[i]
		i++
	}

	// walk up the edge
	m := e.ConnectedNode(n)
	if colors[m] != m {
		// already colored, n is a block root (possibly with cardinality 1)
		return n, n.W
	}

	// check if e crosses an edge with priority
	for _, f := range priority[n.Layer] {
		if crosses(e, f) {
			return n, n.W
		}
	}
	priority[n.Layer] = append(priority[n.Layer], e)

	root, rootw := setColor(m, colors, roots, priority)
	// set m's color to n's, so that eventually each block has the color of its farthest sink node
	colors[m] = n
	roots[n] = root
	return root, max(n.W, rootw)
}

// one run of this routine is O(2n); by placing the recursive call at the end behind a boolean flag, it runs again only once for
// all remaining overlaps. Therefore it becomes O(2*(1+k)*n) where k is the number of times any overlap is found.
func placeBlock(g *graph.DGraph, layerMaxLen int, spacing float64, blockmax, blockwidth, xcoord graph.NodeFloatMap, roots graph.NodeMap) {
	for _, n := range g.Nodes {
		x := blockmax[roots[n]]
		xcoord[n] = max(x, x+(blockwidth[roots[n]]-n.W)/2)
	}

	shift := false
	for k := 0; k < layerMaxLen; k++ {
		for _, l := range g.Layers {
			switch {
			case k >= l.Len():
				continue
			case k == l.Len()-1 && k > 0:
				// consider previous and current nodes
				prv, cur := l.Nodes[k-1], l.Nodes[k]
				// shift if there is an overlap
				if xcoord[cur] < xcoord[prv]+blockwidth[roots[prv]]+spacing {
					xcoord[cur] = xcoord[prv] + blockwidth[roots[prv]] + spacing
					shift = true
					blockmax[roots[cur]] = max(blockmax[roots[cur]], xcoord[cur])
				}
			case k < l.Len()-1:
				// consider current and successive nodes
				cur, suc := l.Nodes[k], l.Nodes[k+1]
				// shift if there is an overlap
				if xcoord[cur] > xcoord[suc] {
					xcoord[suc] = xcoord[cur] + blockwidth[roots[cur]] + spacing
					shift = true
					blockmax[roots[l.Nodes[k+1]]] = max(blockmax[roots[l.Nodes[k+1]]], xcoord[l.Nodes[k+1]])
				}
			}
		}
	}
	if shift {
		placeBlock(g, layerMaxLen, spacing, blockmax, blockwidth, xcoord, roots)
	}
}

func crosses(e, f *graph.Edge) bool {
	if !(e.From.Layer == f.From.Layer && e.To.Layer == f.To.Layer) {
		return false
	}
	etop, ebtm := e.From, e.To
	ftop, fbtm := f.From, f.To
	return (etop.LayerPos < ftop.LayerPos && ebtm.LayerPos > fbtm.LayerPos) ||
		(etop.LayerPos > ftop.LayerPos && ebtm.LayerPos < fbtm.LayerPos)
}

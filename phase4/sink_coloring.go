package phase4

import (
	"github.com/nulab/autog/graph"
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

	// paint nodes, and remember the maximum same-color block width, O(n)
	blockwidth := graph.NodeFloatMap{}
	iter := layersIterator(g, top)
	for layer := iter(); layer != nil; layer = iter() {
		for _, n := range layer.Nodes {
			_, w := setColor(n, colors, roots)
			blockwidth[roots[n]] = max(blockwidth[roots[n]], w)
		}
	}

	// init coordinates by packing nodes to the left, O(n)
	xcoord := graph.NodeFloatMap{}
	iter = layersIterator(g, bottom)
	for layer := iter(); layer != nil; layer = iter() {
		x := 0.0
		for _, n := range layer.Nodes {
			xcoord[n] = x
			x += blockwidth[roots[n]] + params.NodeMargin + params.NodeSpacing
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

	spacing := params.NodeSpacing + params.NodeMargin
	placeBlock(g, lmax, spacing, blockmax, blockwidth, xcoord, roots)

	for _, l := range g.Layers {
		for _, n := range l.Nodes {
			n.X = xcoord[n]
			l.H = max(l.H, n.H)
		}
	}
}

func setColor(n *graph.Node, colors graph.NodeMap, roots graph.NodeMap) (*graph.Node, float64) {
	if colors[n] != n || len(n.In) == 0 {
		return n, n.W
	}

	// candidate edge; this assumes the edge is viable, i.e. doesn't self-loop and is not flat
	mid := len(n.In) / 2
	e := n.In[mid]

	for _, f := range n.In {
		// prefer edges connecting to virtual nodes
		if f.ConnectedNode(n).IsVirtual {
			e = f
		}
	}

	if e == n.In[mid] && len(n.In)%2 == 0 {
		// todo: shift n right
	}
	m := e.ConnectedNode(n)
	if colors[m] != m {
		return n, n.W
	}
	root, rootw := setColor(m, colors, roots)
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

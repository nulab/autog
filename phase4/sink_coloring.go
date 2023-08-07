package phase4

import (
	"github.com/nulab/autog/graph"
)

// todo: add documentation
func execSinkColoring(g *graph.DGraph, params graph.Params) {
	colors := graph.NodeMap{}
	roots := graph.NodeMap{}
	for _, n := range g.Nodes {
		colors[n] = n
		roots[n] = n
	}

	iter := layersIterator(g, top)
	for layer := iter(); layer != nil; layer = iter() {
		for _, n := range layer.Nodes {
			setColor(n, colors, roots)
		}
	}

	xcoord := graph.NodeFloatMap{}
	xinit := graph.NodeSet{}
	blockwidth := graph.NodeFloatMap{}

	placeBlock := func(n *graph.Node, x float64) {
		for {
			xcoord[n] = x
			xinit[n] = true
			blockwidth[n] = max(blockwidth[n], n.W)
			if n != colors[n] {
				n = colors[n]
			} else {
				break
			}
		}
	}

	iter = layersIterator(g, bottom)
	for layer := iter(); layer != nil; layer = iter() {
		x := 0.0
		for _, n := range layer.Nodes {
			x = max(x, xcoord[n])
			if !xinit[n] {
				// apply x coord
				placeBlock(roots[n], x)
			} else {
				if xcoord[n] >= x {
					// do nothing
				} else {
					q := []*graph.Node{n}
					shift := x
					for len(q) > 0 {
						n, q = q[0], q[1:]
						placeBlock(roots[n], shift)
						shift += blockwidth[n] + params.NodeMargin + params.NodeSpacing
						r := roots[n]
						if r.LayerPos < g.Layers[r.Layer].Len()-1 {
							q = append(q, g.Layers[r.Layer].Nodes[r.LayerPos+1])
						}
					}
				}
			}
			x += blockwidth[n] + params.NodeMargin + params.NodeSpacing
		}
	}

	for _, l := range g.Layers {
		for _, n := range l.Nodes {
			n.X = xcoord[n]
			l.H = max(l.H, n.H)
		}
	}
}

func setColor(n *graph.Node, colors graph.NodeMap, roots graph.NodeMap) *graph.Node {
	if colors[n] != n || len(n.In) == 0 {
		return n
	}

	// candidate edge
	mid := len(n.In) / 2
	e := n.In[mid]

	// i := 1
	// j := 0
	// for e.SelfLoops() || e.IsFlat() {
	// 	if j%2 == 0 {
	// 		e = n.In[mid-i]
	// 	} else {
	// 		e = n.In[mid+i]
	// 		i++
	// 	}
	// 	j++
	// }

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
		return n
	}
	root := setColor(m, colors, roots)
	colors[m] = n
	roots[n] = root
	return root
}

func setX(n *graph.Node, colors graph.NodeMap, shift graph.NodeFloatMap, xcoord graph.NodeFloatMap, xinit graph.NodeSet) {
	if xinit[n] || colors[n] == n {
		return
	}
	xcoord[n] = shift[n] // defaults to 0
	xinit[n] = true
	setX(colors[n], colors, shift, xcoord, xinit)
}

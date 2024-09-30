package phase5

import (
	"slices"

	"github.com/nulab/autog/internal/geom"
	"github.com/nulab/autog/internal/graph"
	imonitor "github.com/nulab/autog/internal/monitor"
)

func execSplines(g *graph.DGraph, routes []routableEdge) {
	for _, e := range routes {
		imonitor.Log("spline", e)

		rects := buildRects(g, e)

		for _, r := range rects {
			imonitor.Log("rect", r)
		}

		start := geom.P{e.From.X + e.From.W/2, e.From.Y + e.From.H}
		end := geom.P{e.To.X + e.To.W/2, e.To.Y}

		imonitor.Log("shortest-start", start)
		imonitor.Log("shortest-end", end)
		for _, n := range e.ns {
			imonitor.Log("route-node", n)
		}

		path := geom.Shortest(start, end, rects)

		poly := geom.MergeRects(rects)
		ctrls := geom.FitSpline(path, geom.P{}, geom.P{}, poly.Sides())

		e.Points = make([][2]float64, 0, len(ctrls)*4)
		for _, c := range slices.Backward(ctrls) {
			s := c.Float64Slice()
			e.Points = append(e.Points, [][2]float64{s[3], s[2], s[1], s[0]}...)
		}
	}
}

func buildRects(g *graph.DGraph, r routableEdge) (rects []geom.Rect) {
	for i := 1; i < len(r.ns); i++ {
		top, btm := r.ns[i-1], r.ns[i]
		switch {

		case !top.IsVirtual && !btm.IsVirtual:
			// add one rectangle that spans from the leftmost point to the rightmost point of the two nodes
			r := geom.Rect{
				TL: geom.P{min(top.X, btm.X), top.Y + top.H},
				BR: geom.P{max(top.X+top.W, btm.X+btm.W), btm.Y},
			}
			rects = append(rects, r)

		case btm.IsVirtual:
			// add one rectangle that spans the entire space between the top and bottom layers
			// and one that spans the space around the virtual node
			tl := g.Layers[top.Layer]
			bl := g.Layers[btm.Layer]
			rects = append(rects, rectBetweenLayers(tl, bl))
			rects = append(rects, rectVirtualNode(btm, bl))

		case top.IsVirtual:
			tl := g.Layers[top.Layer]
			bl := g.Layers[btm.Layer]
			rects = append(rects, rectBetweenLayers(tl, bl))
		}
	}

	return
}

func rectBetweenLayers(l1, l2 *graph.Layer) geom.Rect {
	h := l1.Head()
	t := l2.Tail()
	return geom.Rect{
		TL: geom.P{h.X, h.Y + h.H},
		BR: geom.P{t.X + t.W, t.Y},
	}
}

func rectVirtualNode(vn *graph.Node, vl *graph.Layer) geom.Rect {
	switch p := vn.LayerPos; {
	case p == 0:
		// this p+1 access is safe: a layer cannot contain only one virtual node
		n := vl.Nodes[p+1]
		return geom.Rect{
			TL: geom.P{vn.X - 10, n.Y},
			BR: geom.P{n.X, n.Y + n.H},
		}

	case p == vl.Len()-1:
		// this p-1 access is safe: a layer cannot contain only one virtual node
		n := vl.Nodes[p-1]
		return geom.Rect{
			TL: geom.P{n.X + n.W, n.Y},
			BR: geom.P{vn.X + 10, n.Y + n.H},
		}

	default:
		n1 := vl.Nodes[p-1]
		n2 := vl.Nodes[p+1]
		return rectBetweenNodes(n1, n2)
	}
}

func rectBetweenNodes(n1, n2 *graph.Node) geom.Rect {
	d := n2.X - (n1.X + n1.W)
	return geom.Rect{
		TL: geom.P{n1.X + n1.W + d/3, n1.Y},
		BR: geom.P{n2.X - d/3, n2.Y + n2.H},
	}
}

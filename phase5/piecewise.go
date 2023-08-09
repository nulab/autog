package phase5

import "github.com/nulab/autog/graph"

func execPieceWiseRouting(g *graph.DGraph) {
	longEdgesTarget := map[*graph.Node]*graph.Edge{}
	longEdgesSource := map[*graph.Node]*graph.Edge{}

	for i := 0; i < len(g.Edges); i++ {
		e := g.Edges[i]

		switch e.Type() {
		case edgeTypeNoneVirtual:
			e.Points = [][2]float64{startPoint(e), endPoint(e)}

		case edgeTypeOneVirtual:
			if e.From.IsVirtual {
				// source is virtual, check if an edge with the same virtual node as target was encountered
				f := longEdgesTarget[e.From]
				if f != nil {
					f.Points = append(f.Points, bendPoint(e.From, g.Layers[e.From.Layer].H), endPoint(e))
					f.To = e.To
					g.Edges.Remove(e)
					i--
				} else {
					e.Points = append(e.Points, endPoint(e))
					longEdgesSource[e.From] = e
				}
			} else {
				// target is virtual, check if an edge with the same virtual node as source was encountered
				f := longEdgesSource[e.To]
				if f != nil {
					f.Points = append([][2]float64{startPoint(e), bendPoint(e.To, g.Layers[e.To.Layer].H)}, f.Points...)
					f.From = e.From
					g.Edges.Remove(e)
					i--
				} else {
					e.Points = append(e.Points, startPoint(e))
					longEdgesTarget[e.To] = e
				}
			}

		case edgeTypeBothVirtual:
			f := longEdgesTarget[e.From]
			if f != nil {
				f.Points = append(f.Points, bendPoint(e.From, g.Layers[e.From.Layer].H))
				f.To = e.To
				longEdgesTarget[e.To] = f
			} else {
				longEdgesTarget[e.To] = e
			}
			f = longEdgesSource[e.To]
			if f != nil {
				f.From = e.From
				longEdgesSource[e.From] = f
			} else {
				longEdgesSource[e.From] = e
			}
			g.Edges.Remove(e)
			i--
		}
	}
}

func bendPoint(n *graph.Node, layerHeight float64) [2]float64 {
	if !n.IsVirtual {
		panic("routing: bend point on non-virtual node")
	}
	var x, y float64
	x = n.X + n.W/2
	y = n.Y + layerHeight/2
	return [2]float64{x, y}
}
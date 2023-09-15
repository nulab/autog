package phase5

import "github.com/nulab/autog/internal/graph"

func execStraightRouting(g *graph.DGraph) {
	longEdgesTarget := map[*graph.Node]*graph.Edge{}
	longEdgesSource := map[*graph.Node]*graph.Edge{}

	for i := 0; i < len(g.Edges); i++ {
		e := g.Edges[i]

		switch e.Type() {
		case edgeTypeNoneVirtual:
			if e.IsFlat() {
				e.Points = [][2]float64{flatStartPoint(e), flatEndPoint(e)}
			} else {
				e.Points = [][2]float64{startPoint(e), endPoint(e)}
			}

		case edgeTypeOneVirtual:
			if e.From.IsVirtual {
				// source is virtual, check if an edge with the same virtual node as target was encountered
				f := longEdgesTarget[e.From]
				if f != nil {
					f.Points = append(f.Points, endPoint(e)) // e's endpoint becomes f's endpoint
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
					f.Points = append([][2]float64{startPoint(e)}, f.Points...) // e's startpoint becomes f's startpoint
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
				// f.Points = append(f.Points, endPoint(e)) // e's endpoint becomes f's endpoint
				f.To = e.To
				longEdgesTarget[e.To] = f
			} else {
				// e.Points = append(e.Points, endPoint(e))
				longEdgesTarget[e.To] = e
			}
			f = longEdgesSource[e.To]
			if f != nil {
				// f.Points = append([][2]float64{startPoint(e)}, f.Points...) // e's startpoint becomes f's startpoint
				f.From = e.From
				longEdgesSource[e.From] = f
			} else {
				// e.Points = append(e.Points, startPoint(e))
				longEdgesSource[e.From] = e
			}
			g.Edges.Remove(e)
			i--
		}
	}
}

func flatStartPoint(e *graph.Edge) [2]float64 {
	var x, y float64
	if e.From.LayerPos < e.To.LayerPos {
		// middle of right side
		x = e.From.X + e.From.W
		y = e.From.Y + e.From.H/2
	} else {
		// middle of left side
		x = e.From.X
		y = e.From.Y + e.From.H/2
	}
	return [2]float64{x, y}
}

func flatEndPoint(e *graph.Edge) [2]float64 {
	var x, y float64
	if e.From.LayerPos < e.To.LayerPos {
		// middle of left side
		x = e.To.X
		y = e.To.Y + e.To.H/2
	} else {
		// middle of right side
		x = e.To.X + e.To.W
		y = e.To.Y + e.To.H/2
	}
	return [2]float64{x, y}
}

func startPoint(e *graph.Edge) [2]float64 {
	var x, y float64
	if e.From.Layer < e.To.Layer {
		// middle of lower side
		x = e.From.X + e.From.W/2
		y = e.From.Y + e.From.H
	} else {
		// middle of upper side
		x = e.From.X + e.From.W/2
		y = e.From.Y
	}
	return [2]float64{x, y}
}

func endPoint(e *graph.Edge) [2]float64 {
	var x, y float64
	if e.From.Layer < e.To.Layer {
		// middle of upper side
		x = e.To.X + e.To.W/2
		y = e.To.Y
	} else {
		// middle of lower side
		x = e.To.X + e.To.W/2
		y = e.To.Y + e.To.H
	}
	return [2]float64{x, y}
}

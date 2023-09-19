package phase5

import (
	"slices"

	"github.com/nulab/autog/internal/graph"
	"github.com/nulab/autog/internal/num"
)

func execOrthoRouting(g *graph.DGraph, params graph.Params) {
	longEdgesTarget := map[*graph.Node]*graph.Edge{}
	longEdgesSource := map[*graph.Node]*graph.Edge{}

	halfLayerSpacing := params.LayerSpacing / 2

	for i := 0; i < len(g.Edges); i++ {
		e := g.Edges[i]

		switch e.Type() {
		case edgeTypeNoneVirtual:
			if e.IsFlat() {
				if num.Abs(e.From.LayerPos-e.To.LayerPos) > 1 {
					e.Points = flatNonConsecutive(e, g.Layers[e.From.Layer].H)
				} else {
					e.Points = [][2]float64{flatStartPoint(e), flatEndPoint(e)}
				}
			} else {
				e.Points = orthoPoints(e, halfLayerSpacing)
			}

		case edgeTypeOneVirtual:
			ps := orthoPoints(e, halfLayerSpacing)
			if e.From.IsVirtual {
				// source is virtual, check if an edge with the same virtual node as target was encountered
				f := longEdgesTarget[e.From]
				if f != nil {
					f.Points = append(f.Points, ps...)
					f.To = e.To
					g.Edges.Remove(e)
					i--
				} else {
					e.Points = append(e.Points, ps...)
					longEdgesSource[e.From] = e
				}
			} else {
				// target is virtual, check if an edge with the same virtual node as source was encountered
				f := longEdgesSource[e.To]
				if f != nil {
					f.Points = append(ps, f.Points...)
					f.From = e.From
					g.Edges.Remove(e)
					i--
				} else {
					e.Points = append(e.Points, ps...)
					longEdgesTarget[e.To] = e
				}
			}

		case edgeTypeBothVirtual:
			f := longEdgesTarget[e.From]
			if f != nil {
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

func orthoPoints(e *graph.Edge, halfLayerSpacing float64) [][2]float64 {
	if isVerticallyAligned(e) {
		return [][2]float64{startPoint(e.From), endPoint(e.To)}
	}

	y := e.From.Y + e.From.H
	if e.From.IsVirtual {
		y = e.To.Y - halfLayerSpacing*2
	}
	ps := make([][2]float64, 0, 5)
	ps = append(ps, [2]float64{e.From.X + e.From.W/2, y})
	ps = append(ps, [2]float64{e.From.X + e.From.W/2, y + halfLayerSpacing})
	ps = append(ps, [2]float64{e.To.X + e.To.W/2, e.To.Y - halfLayerSpacing})
	ps = append(ps, [2]float64{e.To.X + e.To.W/2, e.To.Y})
	if e.From.Layer > e.To.Layer {
		slices.Reverse(ps)
	}
	return ps
}

func isVerticallyAligned(e *graph.Edge) bool {
	return e.From.X+e.From.W/2 == e.To.X+e.To.W/2
}

func flatStartPoint(e *graph.Edge) [2]float64 {
	var x, y float64
	if e.From.LayerPos < e.To.LayerPos {

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

	} else {
		// middle of right side
		x = e.To.X + e.To.W
		y = e.To.Y + e.To.H/2
	}
	return [2]float64{x, y}
}

package phase5

import (
	"github.com/nulab/autog/internal/graph"
)

func execOrthoRouting(g *graph.DGraph, routes []routableEdge, params graph.Params) {
	halfLayerSpacing := params.LayerSpacing / 2
	for _, r := range routes {
		layerh := g.Layers[r.From.Layer].H
		if r.IsFlat() {
			flatPolyline(r, layerh)
			continue
		}
		if isVerticallyAligned(r.From, r.To) {
			r.Points = straight(r.ns[0], r.ns[len(r.ns)-1])
			continue
		}

		for i := 1; i < len(r.ns); i++ {
			sp := startPoint(r.ns[i-1])
			// virtual nodes have 0 size; another solution here is to consider the layer Y instead of the node Y
			if r.ns[i-1].IsVirtual {
				sp[1] += layerh
			}
			r.Points = append(r.Points, sp)
			r.Points = append(r.Points, [2]float64{sp[0], sp[1] + halfLayerSpacing})

			ep := endPoint(r.ns[i])
			r.Points = append(r.Points, [2]float64{ep[0], ep[1] - halfLayerSpacing})
			r.Points = append(r.Points, ep)
		}
	}
}

func isVerticallyAligned(from, to *graph.Node) bool {
	return from.X+from.W/2 == to.X+to.W/2
}

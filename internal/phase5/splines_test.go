package phase5

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nulab/autog/internal/geom"
	"github.com/nulab/autog/internal/graph"
)

func TestRectBetweenLayers(t *testing.T) {
	node := func(id string, x, y, w, h float64) *graph.Node {
		return &graph.Node{ID: id, Size: graph.Size{X: x, Y: y, W: w, H: h}}
	}
	layer := func(h float64, nodes ...*graph.Node) *graph.Layer {
		return &graph.Layer{Nodes: nodes, Size: graph.Size{H: h}}
	}

	testCases := []struct {
		name   string
		l1, l2 *graph.Layer
		want   geom.Rect
	}{
		{
			// regression: the right bound must account for the upper layer's tail
			name: "upper layer extends further right",
			l1:   layer(40, node("A", 0, 0, 50, 40), node("B", 200, 0, 100, 40)),
			l2:   layer(40, node("C", 20, 190, 60, 40), node("D", 120, 190, 60, 40)),
			want: geom.Rect{TL: geom.P{0, 40}, BR: geom.P{300, 190}},
		},
		{
			name: "lower layer extends further left and right",
			l1:   layer(40, node("A", 50, 0, 40, 40)),
			l2:   layer(40, node("C", 0, 190, 40, 40), node("D", 200, 190, 80, 40)),
			want: geom.Rect{TL: geom.P{0, 40}, BR: geom.P{280, 190}},
		},
		{
			name: "upper layer extends further left, lower layer further right",
			l1:   layer(40, node("A", 0, 0, 40, 40), node("B", 100, 0, 40, 40)),
			l2:   layer(40, node("C", 60, 190, 40, 40), node("D", 180, 190, 40, 40)),
			want: geom.Rect{TL: geom.P{0, 40}, BR: geom.P{220, 190}},
		},
		{
			name: "upper layer taller than its head node",
			l1:   layer(60, node("A", 0, 0, 40, 40), node("B", 100, 0, 40, 60)),
			l2:   layer(40, node("C", 0, 210, 40, 40)),
			want: geom.Rect{TL: geom.P{0, 60}, BR: geom.P{140, 210}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, rectBetweenLayers(tc.l1, tc.l2))
		})
	}
}

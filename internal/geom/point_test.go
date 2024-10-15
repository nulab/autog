package geom

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint(t *testing.T) {
	t.Run("svg", func(t *testing.T) {
		p := P{10.8, 50.1}
		assert.Equal(t, `<circle r="4" cx="10.80" cy="50.10" fill="black"/>`, p.SVG())
	})

	t.Run("addp", func(t *testing.T) {
		a := P{1, 2}
		b := P{3, 4}
		assert.Equal(t, P{4, 6}, addp(a, b))
	})

	t.Run("subp", func(t *testing.T) {
		a := P{5, 6}
		b := P{3, 4}
		assert.Equal(t, P{2, 2}, subp(a, b))
	})

	t.Run("scalep", func(t *testing.T) {
		p := P{2, 3}
		assert.Equal(t, P{4, 6}, scalep(p, 2.0))
	})

	t.Run("dotp", func(t *testing.T) {
		a := P{1, 2}
		b := P{3, 4}
		assert.Equal(t, 11.0, dotp(a, b))
	})

	t.Run("distp", func(t *testing.T) {
		a := P{0, 0}
		b := P{3, 4}
		assert.Equal(t, 5.0, distp(a, b))
	})

	t.Run("sqdistp", func(t *testing.T) {
		a := P{0, 0}
		b := P{3, 4}
		assert.Equal(t, 25.0, sqdistp(a, b))
	})

	t.Run("norm", func(t *testing.T) {
		a := P{3, 4}
		got := norm(a)
		want := P{0.6, 0.8}
		assert.InDelta(t, want.X, got.X, 1e-9)
		assert.InDelta(t, want.Y, got.Y, 1e-9)
	})

	t.Run("rotatep", func(t *testing.T) {
		p := P{1, 0}
		got := rotatep(p, math.Pi/2)
		want := P{0, 1}
		assert.InDelta(t, want.X, got.X, 1e-9)
		assert.InDelta(t, want.Y, got.Y, 1e-9)
	})
}

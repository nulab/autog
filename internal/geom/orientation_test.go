package geom

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrientation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		a, b, c := P{0, 0}, P{20, 20}, P{40, 20}
		assert.Equal(t, ccw, orientation(a, b, c))
		assert.Equal(t, cw, orientation(a, c, b))
		assert.Equal(t, cw, orientation(c, b, a))
		assert.Equal(t, ccw, orientation(c, a, b))
	})

	t.Run("success with small margin", func(t *testing.T) {
		a, b, c := P{12.0, 4.78}, P{12.0000000001, 14.78}, P{12.0, 20}
		assert.Equal(t, cw, orientation(a, b, c))
		assert.Equal(t, ccw, orientation(a, c, b))
	})

	t.Run("collinear single point", func(t *testing.T) {
		a := P{math.Pi, math.Cos(78.12)}
		b, c := a, a
		assert.Equal(t, cln, orientation(a, b, c))
	})

	t.Run("collinear two identical", func(t *testing.T) {
		a := P{45.787232, 34.9829283}
		c := P{12.8934, 65.9232}
		b := a
		assert.Equal(t, cln, orientation(a, b, c))
	})

	t.Run("collinear origin", func(t *testing.T) {
		a, b, c := P{}, P{}, P{}
		assert.Equal(t, cln, orientation(a, b, c))
	})

	t.Run("collinear parallel to x axis", func(t *testing.T) {
		a, b, c := P{-8934.12, 50.566}, P{0, 50.566}, P{math.Sqrt(23), 50.566}
		assert.Equal(t, cln, orientation(a, b, c))
	})

	t.Run("collinear parallel to y axis", func(t *testing.T) {
		a, b, c := P{20.001, 50}, P{20.001, 75.346}, P{20.001, -0.00000001}
		assert.Equal(t, cln, orientation(a, b, c))
	})

	t.Run("collinear with pos slope", func(t *testing.T) {
		f := func(x float64) float64 {
			return 4*x/5.0 + 7.58
		}
		x1, x2, x3 := 45.68, 0.012, -746.1
		a, b, c := P{x1, f(x1)}, P{x2, f(x2)}, P{x3, f(x3)}
		assert.Equal(t, cln, orientation(a, b, c))
	})
}

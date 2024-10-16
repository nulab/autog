package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeSpline(t *testing.T) {
	const delta = 1e-2

	t.Run("vertically aligned", func(t *testing.T) {
		a := P{20, 30}
		b := P{20, 50}
		s := MakeSpline(a, b)
		assert.True(t, s.p0 == s.p1)
		assert.True(t, s.p2 == s.p3)
	})

	t.Run("horizontal segment", func(t *testing.T) {
		a := P{20, 10}
		b := P{50, 10}
		s := MakeSpline(a, b)
		assert.Equal(t, s.p0, a)
		assert.InDelta(t, 25.7, s.p1.X, delta)
		assert.InDelta(t, 8.14, s.p1.Y, delta)
		assert.InDelta(t, 44.3, s.p2.X, delta)
		assert.InDelta(t, 8.14, s.p2.Y, delta)
		assert.Equal(t, s.p3, b)
	})

	t.Run("negative slope", func(t *testing.T) {
		a := P{50, 20}
		b := P{20, 10}
		s := MakeSpline(a, b)
		assert.Equal(t, s.p0, a)
		assert.InDelta(t, 44.91, s.p1.X, delta)
		assert.InDelta(t, 16.24, s.p1.Y, delta)
		assert.InDelta(t, 26.32, s.p2.X, delta)
		assert.InDelta(t, 10.04, s.p2.Y, delta)
		assert.Equal(t, s.p3, b)
	})

	t.Run("positive slope", func(t *testing.T) {
		a := P{20, 10}
		b := P{50, 20}
		s := MakeSpline(a, b)
		assert.Equal(t, s.p0, a)
		assert.InDelta(t, 26.32, s.p1.X, delta)
		assert.InDelta(t, 10.04, s.p1.Y, delta)
		assert.InDelta(t, 44.91, s.p2.X, delta)
		assert.InDelta(t, 16.24, s.p2.Y, delta)
		assert.Equal(t, s.p3, b)
	})
}

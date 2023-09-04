package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTriangulate(t *testing.T) {
	t.Run("ordered points", func(t *testing.T) {
		p := P{5.5, 7.9}
		q := P{5.6, 7.9}
		assert.Equal(t, q, rightmost(p, q))
		assert.Equal(t, p, leftmost(p, q))
		a, b := left2right(p, q)
		assert.Equal(t, a, p)
		assert.Equal(t, b, q)

		// equal x'es
		p = P{5.5, 7.9}
		q = P{5.5, 8.0}
		assert.Equal(t, p, rightmost(p, q))
		assert.Equal(t, q, leftmost(p, q))
	})
}

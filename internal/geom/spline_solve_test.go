package geom

import (
	"math"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolve(t *testing.T) {
	t.Run("solve cubic three roots", func(t *testing.T) {
		a, b, c, d := -1.0, 7.0, -4.0, -10.0
		roots := solve3([]float64{d, c, b, a})
		assert.Len(t, roots, 3)
		sort.Float64s(roots)
		assert.InEpsilon(t, -0.900052, roots[0], epsilon2)
		assert.InEpsilon(t, 1.830534, roots[1], epsilon2)
		assert.InEpsilon(t, 6.069517, roots[2], epsilon2)
	})

	t.Run("solve cubic one root", func(t *testing.T) {
		a, b, c, d := 5.0, 0.0, 2.0, -2.0
		roots := solve3([]float64{d, c, b, a})
		assert.Len(t, roots, 1)
		assert.InEpsilon(t, 0.560286, roots[0], epsilon2)
	})

	t.Run("solve quadratic root", func(t *testing.T) {
		a, b, c, d := 0.0, math.E, -5.45, 2.0
		roots := solve3([]float64{d, c, b, a})
		assert.Len(t, roots, 2)
		sort.Float64s(roots)
		assert.InEpsilon(t, 0.4836360, roots[0], epsilon2)
		assert.InEpsilon(t, 1.521306, roots[1], epsilon2)
	})

}

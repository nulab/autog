package phase5

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlg(t *testing.T) {
	assert.EqualValues(t, 5, _endAlg)

	strs := []string{"noop", "straight", "piecewise", "ortho", "splines"}

	for i := Alg(0); i < _endAlg; i++ {
		assert.Equal(t, 5, i.Phase())
		assert.Equal(t, strs[i], i.String())
	}
}

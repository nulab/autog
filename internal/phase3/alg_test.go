package phase3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlg(t *testing.T) {
	assert.EqualValues(t, 2, _endAlg)

	strs := []string{"noop", "gvdot"}

	for i := Alg(0); i < _endAlg; i++ {
		assert.Equal(t, 3, i.Phase())
		assert.Equal(t, strs[i], i.String())
	}
	assert.Equal(t, "<invalid>", _endAlg.String())
}

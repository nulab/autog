package phase1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlg(t *testing.T) {
	assert.EqualValues(t, 2, _endAlg)

	strs := []string{"greedy", "dfs"}

	for i := Alg(0); i < _endAlg; i++ {
		assert.Equal(t, 1, i.Phase())
		assert.Equal(t, strs[i], i.String())
	}
}

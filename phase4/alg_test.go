package phase4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlg(t *testing.T) {
	assert.EqualValues(t, 6, _endAlg)

	strs := []string{"noop", "vertical", "b&k", "ns", "sinkcoloring", "packright"}

	for i := Alg(0); i < _endAlg; i++ {
		assert.Equal(t, 4, i.Phase())
		assert.Equal(t, strs[i], i.String())
	}
}

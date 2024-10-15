package collectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMat(t *testing.T) {
	var m Mat[int]

	m = NewMat[int](0)
	assert.Len(t, m, 0)

	m = NewMat[int](3)
	assert.Len(t, m, 3)

	m = NewMat[int](-1)
	assert.Len(t, m, 0)
}

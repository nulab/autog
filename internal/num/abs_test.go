package num

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	assert.Equal(t, 192, Abs(-192))
	assert.Equal(t, 56, Abs(56))
	assert.Equal(t, 0, Abs(0))
}

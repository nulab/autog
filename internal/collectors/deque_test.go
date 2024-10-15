package collectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeque(t *testing.T) {
	deq := NewDeque[string](10)
	// w2 w1 u p v w3
	deq.PushFront("p")
	deq.PushFront("u")
	deq.PushBack("v")
	deq.PushFront("w1")
	deq.PushFront("w2")
	deq.PushBack("w3")
	assert.Equal(t, 6, deq.Len())
	assert.Equal(t, 6, deq.Front()) // 10-4
	assert.Equal(t, 11, deq.Back()) // 9+2
	assert.Equal(t, "v", deq.data[10])
	assert.Equal(t, "w2", deq.PopFront())
	assert.Equal(t, "w3", deq.PopBack())
	assert.Equal(t, "w1", deq.PeekFront(1))
	assert.Equal(t, "v", deq.PeekBack(1))
}

package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPairingSimple(t *testing.T) {
	queue := NewPairing()
	queue.Push("first", 1)
	queue.Push("second", 2)
	queue.Push("third", 3)
	assert.Equal(t, "first", queue.Pop())
	assert.Equal(t, "second", queue.Pop())
	assert.Equal(t, "third", queue.Pop())
	// No more items
	assert.Equal(t, nil, queue.Pop())
}

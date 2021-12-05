package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinkedSimple(t *testing.T) {
	queue := NewLinked()
	queue.Push("second", 2)
	queue.Push("first", 0)
	queue.Push("third", 5)
	assert.Equal(t, "first", queue.Pop())
	assert.Equal(t, "second", queue.Pop())
	assert.Equal(t, "third", queue.Pop())
	// No more items
	assert.Equal(t, nil, queue.Pop())
}

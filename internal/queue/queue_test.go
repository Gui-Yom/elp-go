package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testQueue(t *testing.T, queue PriorityQueue) {
	queue.Push("second", 2)
	queue.Push("first", 0)
	queue.Push("third", 5)
	assert.Equal(t, "first", queue.Pop())
	assert.Equal(t, "second", queue.Pop())
	assert.Equal(t, "third", queue.Pop())
	// No more items
	assert.Equal(t, nil, queue.Pop())
}

func TestQueueLinked(t *testing.T) {
	testQueue(t, NewLinked())
}

func TestQueuePairing(t *testing.T) {
	testQueue(t, NewPairing())
}

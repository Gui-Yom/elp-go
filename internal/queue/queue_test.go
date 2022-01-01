package queue

import (
	"elp-go/internal/world"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testQueue(t *testing.T, queue PriorityQueue) {
	queue.Push(world.Pos(2, 2), 2)
	queue.Push(world.Pos(0, 0), 0)
	queue.Push(world.Pos(5, 5), 5)
	queue.Push(world.Pos(1, 1), 1)
	queue.Push(world.Pos(4, 4), 4)
	assert.Equal(t, world.Pos(0, 0), queue.Pop())
	assert.Equal(t, world.Pos(1, 1), queue.Pop())
	assert.Equal(t, world.Pos(2, 2), queue.Pop())
	assert.Equal(t, world.Pos(4, 4), queue.Pop())
	assert.Equal(t, world.Pos(5, 5), queue.Pop())
	// No more items
	assert.Panics(t, func() { queue.Pop() })
}

func TestQueueLinked(t *testing.T) {
	testQueue(t, NewLinked())
}

func TestQueuePairing(t *testing.T) {
	testQueue(t, NewPairing())
}

func TestQueueHeap(t *testing.T) {
	testQueue(t, NewHeap())
}

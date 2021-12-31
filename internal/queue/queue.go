package queue

import (
	"container/list"
	"elp-go/internal/world"
)

// PriorityQueue A min-priority queue.
type PriorityQueue interface {
	// Push Insert an item with the specified priority
	Push(item world.Position, priority float32)
	// Pop Removes an item, returns nil if empty
	Pop() world.Position
	// Empty returns true if there are no more items
	Empty() bool
}

// NewPairing Creates a new priority queue based on a pairing heap
func NewPairing() PriorityQueue {
	return &pairingQueue{}
}

// NewLinked Creates a new priority queue based a sorted linked list insert.
func NewLinked() PriorityQueue {
	return linkedQueue{items: list.New()}
}

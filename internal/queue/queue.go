package queue

import "container/list"

// PriorityQueue A min-priority queue.
type PriorityQueue interface {
	// Push Insert an item with the specified priority
	Push(item interface{}, priority float32)
	// Pop Removes an item, returns nil if empty
	Pop() interface{}
	// Empty returns true if there are no more items
	Empty() bool
}

// NewPairing Creates a new priority queue based on a pairing heap
func NewPairing() PriorityQueue {
	return pairingQueue{}
}

// NewNaive Creates a new priority queue based a naive sorted list insert.
func NewNaive() PriorityQueue {
	return naiveQueue{items: list.New()}
}

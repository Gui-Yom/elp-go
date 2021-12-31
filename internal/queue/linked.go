package queue

import (
	"container/list"
	"elp-go/internal/world"
)

// Force implementation, thanks Go
var _ PriorityQueue = (*linkedQueue)(nil)

// linkedQueue Priority queue based on a linked list.
// The list implementation is provided by the Go stdlib and is actually a doubly linked list.
// We might get some more juice with a simple linked list.
type linkedQueue struct {
	items *list.List
}

// lqnode the linked list node, holding our item with its priority
// NOTE: our implementation is generic over the item values (interface{}), this costs one more allocation each time because the item gets boxed.
// We should really specialize the type because we only use those priority queues to store pathfinding.Position items.
type lqnode struct {
	item     world.Position
	priority float32
}

func (n linkedQueue) Push(item world.Position, priority float32) {
	element := lqnode{item: item, priority: priority}
	if n.items.Len() == 0 {
		n.items.PushFront(element)
	} else {
		for curr := n.items.Front(); curr != nil; curr = curr.Next() {
			// We insert before the first element with the same or bigger priority
			// This does not preserve insertion order and thus our implementation is not stable (do we care tho ?)
			// But if multiple items have the same priority (equal distance to the goal) this might save a few iterations
			if priority <= curr.Value.(lqnode).priority {
				n.items.InsertBefore(element, curr)
				return
			}
		}
		n.items.PushBack(element)
	}
}

func (n linkedQueue) Pop() world.Position {
	elem := n.items.Front()
	if elem == nil {
		panic("Tried to Pop() with no items")
		return world.Position{}
	} else {
		n.items.Remove(elem)
		return elem.Value.(lqnode).item
	}
}

func (n linkedQueue) Empty() bool {
	return n.items.Len() == 0
}

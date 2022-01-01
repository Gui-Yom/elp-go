package queue

import (
	"elp-go/internal/world"
)

// Force implementation, thanks Go
var _ PriorityQueue = (*linkedQueue)(nil)

type linkedQueue struct {
	root *lqnode
}

type lqnode struct {
	item world.Position
	cost float64
	next *lqnode
}

func (list *linkedQueue) Empty() bool {
	return list.root == nil
}

func (list *linkedQueue) Push(pos world.Position, cost float64) {
	newNode := lqnode{
		item: pos,
		cost: cost,
	}
	if list.root == nil { // List is empty
		list.root = &newNode
		return
	}
	if cost <= list.root.cost { // List has 1 item, this removes the branch in the loop
		newNode.next = list.root
		list.root = &newNode
		return
	}
	var pred = list.root
	for curr := pred.next; curr != nil; curr = curr.next { // Iterate through values
		// We insert before the first element with the same or bigger priority
		// This does not preserve insertion order and thus our implementation is not stable (do we care tho ?)
		// But if multiple items have the same priority (equal distance to the goal) this might save a few iterations
		if cost <= curr.cost {
			newNode.next = curr
			pred.next = &newNode
			return
		}
		pred = curr
	}
	pred.next = &newNode
}

func (list *linkedQueue) Pop() world.Position {
	if list.root == nil {
		panic("Tried to Pop() with no items")
		return world.Position{}
	} else {
		pos := list.root.item
		list.root = list.root.next
		return pos
	}
}

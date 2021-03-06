package queue

import (
	"elp-go/internal/world"
)

var _ PriorityQueue = (*pairingQueue)(nil)

// pairingQueue a priority queue based on a pairing heap.
// see https://en.wikipedia.org/wiki/Pairing_heap.
type pairingQueue struct {
	root *pqnode
}

type pqnode struct {
	// No generics, thanks Go.
	item world.Position
	// The priority of this item
	priority float64
	children []*pqnode
}

func (q *pairingQueue) Push(item world.Position, priority float64) {
	q.root = merge(q.root, &pqnode{item: item, priority: priority})
}

func (q *pairingQueue) Empty() bool {
	return q.root == nil
}

func (q *pairingQueue) Pop() world.Position {
	if q.root == nil {
		panic("Tried to Pop() with no items")
		return world.Position{}
	}
	item := q.root.item
	q.root = mergeChildren(q.root, q.root.children)
	return item
}

func merge(a *pqnode, b *pqnode) *pqnode {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.priority < b.priority {
		a.children = append([]*pqnode{b}, a.children...)
		return a
	} else {
		b.children = append([]*pqnode{a}, b.children...)
		return b
	}
}

func mergeChildren(root *pqnode, heaps []*pqnode) *pqnode {
	if len(heaps) == 1 {
		root = heaps[0]
		return root
	}
	var merged *pqnode
	for {
		if len(heaps) == 0 {
			break
		}
		if merged == nil {
			merged = merge(heaps[0], heaps[1])
			heaps = heaps[2:]
		} else {
			merged = merge(merged, heaps[0])
			heaps = heaps[1:]
		}
	}
	root = merged

	return root
}

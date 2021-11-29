package internal

// PriorityQueue a priority queue based on a pairing heap.
// see https://en.wikipedia.org/wiki/Pairing_heap.
type PriorityQueue struct {
	root *node
}

type node struct {
	// No generics, thanks Go.
	item interface{}
	// The priority of this item
	priority float32
	children []*node
}

// push insert an item in the queue.
func (q *PriorityQueue) push(item interface{}, priority float32) {
	q.root = merge(q.root, &node{item: item, priority: priority})
}

func (q *PriorityQueue) empty() bool {
	return q.root == nil
}

// pop retrieve an item from the queue
func (q *PriorityQueue) pop() interface{} {
	if q.root == nil {
		return nil
	}
	item := q.root.item
	q.root = mergeChildren(q.root, q.root.children)
	return item
}

func merge(a *node, b *node) *node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.priority < b.priority {
		a.children = append([]*node{b}, a.children...)
		return a
	} else {
		b.children = append([]*node{a}, b.children...)
		return b
	}
}

func mergeChildren(root *node, heaps []*node) *node {
	if len(heaps) == 1 {
		root = heaps[0]
		return root
	}
	var merged *node
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

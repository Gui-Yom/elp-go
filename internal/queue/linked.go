package queue

import "container/list"

var _ PriorityQueue = (*linkedQueue)(nil)

type linkedQueue struct {
	items *list.List
}

type lqnode struct {
	item     interface{}
	priority float32
}

func (n linkedQueue) Push(item interface{}, priority float32) {
	element := lqnode{item: item, priority: priority}
	if n.items.Len() == 0 {
		n.items.PushFront(element)
	} else {
		for curr := n.items.Front(); curr != nil; curr = curr.Next() {
			if priority < curr.Value.(lqnode).priority {
				n.items.InsertBefore(element, curr)
				return
			}
		}
		n.items.PushBack(element)
	}
}

func (n linkedQueue) Pop() interface{} {
	elem := n.items.Front()
	if elem == nil {
		return nil
	} else {
		n.items.Remove(elem)
		return elem.Value.(lqnode).item
	}
}

func (n linkedQueue) Empty() bool {
	return n.items.Len() == 0
}

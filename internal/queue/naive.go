package queue

import "container/list"

var _ PriorityQueue = (*naiveQueue)(nil)

type naiveQueue struct {
	items *list.List
}

type nqnode struct {
	item     interface{}
	priority float32
}

func (n naiveQueue) Push(item interface{}, priority float32) {
	/*
		n.items = append(n.items, pqnode{item: item, priority: priority})
		sort.Slice(n.items, func(i, j int) bool {
			return n.items[i].priority < n.items[j].priority
		})
	*/
	element := nqnode{item: item, priority: priority}
	if n.items.Len() == 0 {
		n.items.PushFront(element)
	} else {
		for curr := n.items.Front(); curr != nil; curr = curr.Next() {
			if priority < curr.Value.(nqnode).priority {
				n.items.InsertBefore(element, curr)
				return
			}
		}
		n.items.PushBack(element)
	}
}

func (n naiveQueue) Pop() interface{} {
	elem := n.items.Front()
	if elem == nil {
		return nil
	} else {
		n.items.Remove(elem)
		return elem.Value.(nqnode).item
	}
}

func (n naiveQueue) Empty() bool {
	return n.items.Len() == 0
}

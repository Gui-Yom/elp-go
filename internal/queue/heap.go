package queue

import "elp-go/internal/world"

type heapQueue []heapNode

var _ PriorityQueue = (*heapQueue)(nil)

type heapNode struct {
	item world.Position
	cost float64
}

// Push insert en item
// O(log n)
func (h *heapQueue) Push(item world.Position, priority float64) {
	n := len(*h)
	*h = append(*h, heapNode{item: item, cost: priority})
	h.up(n - 1)
}

// Pop retrieve an item
// O(log n)
func (h *heapQueue) Pop() world.Position {
	n := len(*h) - 1
	h.swap(0, n)
	h.down(0, n)
	old := *h
	n = len(old)
	node := old[n-1]
	*h = old[0 : n-1]
	return node.item
}

func (h *heapQueue) Empty() bool {
	return len(*h) == 0
}

func (h heapQueue) less(i, j int) bool {
	return h[i].cost < h[j].cost
}

func (h heapQueue) swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h heapQueue) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h heapQueue) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}

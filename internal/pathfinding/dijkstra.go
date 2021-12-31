package pathfinding

import (
	"elp-go/internal/queue"
	"time"
)

type Dijkstra struct {
	Diagonal     bool
	queueBuilder func() queue.PriorityQueue
}

// Implementation is implicit, thanks Go
var _ Pathfinder = (*Dijkstra)(nil)

func (dijk Dijkstra) FindPath(world *World, start Position, goal Position) ([]Position, Stats) {
	startTime := time.Now()

	costs := make(map[Position]float32)
	costs[start] = 0

	frontier := dijk.queueBuilder()
	frontier.Push(start, 0)
	parentChain := make(map[Position]Position)

	var iter uint
	var curr Position

	for !frontier.Empty() {
		curr = frontier.Pop().(Position)

		if curr == goal {
			break
		}

		iter++

		neighbors := world.GetNeighbors(curr, dijk.Diagonal)
		for _, node := range neighbors {
			tileCost := world.GetCost(node)
			newCost := costs[curr] + tileCost
			prevCost, exists := costs[node]
			if !exists || newCost < prevCost {
				costs[node] = newCost
				parentChain[node] = curr
				frontier.Push(node, newCost)
			}
		}
	}
	if curr != goal {
		return nil, Stats{Iterations: iter, Duration: time.Now().Sub(startTime)}
	}
	path := makePath(parentChain, start, goal)
	return path, Stats{Iterations: iter, Duration: time.Now().Sub(startTime), Cost: pathCost(world, path)}
}

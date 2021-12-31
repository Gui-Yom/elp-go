package pathfinding

import (
	"elp-go/internal/queue"
	"math"
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

	// Dijkstra explores everywhere, we model it by a disc pi*r^2
	// We divide by 8 to account for the fact that most of the disc is out of bounds or walls
	presize := int(math.Max(math.Pi*float64(EuclideanSq(start, goal))/8.0, 4))
	//fmt.Printf("presize: %v\n", presize)
	costs := make(map[Position]float32, presize)
	costs[start] = 0

	frontier := dijk.queueBuilder()
	frontier.Push(start, 0)
	parentChain := make(map[Position]Position, presize)

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
	//fmt.Printf("costLen: %v, chainLen: %v\n", len(costs), len(parentChain))
	if curr != goal {
		return nil, Stats{Iterations: iter, Duration: time.Now().Sub(startTime)}
	}
	path := makePath(parentChain, start, goal)
	return path, Stats{Iterations: iter, Duration: time.Now().Sub(startTime), Cost: pathCost(world, path)}
}

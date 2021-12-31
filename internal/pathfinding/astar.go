package pathfinding

import (
	"elp-go/internal/queue"
	"math"
	"time"
)

type Astar struct {
	Diagonal     bool
	Heuristic    Heuristic
	queueBuilder func() queue.PriorityQueue
}

var _ Pathfinder = (*Astar)(nil)

func (astar Astar) FindPath(world *World, start Position, goal Position) ([]Position, Stats) {
	startTime := time.Now()

	// Astar explores in the direction of the goal, we model it by a segment
	// We multiply that number to account for the fact there is walls
	presize := int(math.Max(2*math.Sqrt2*float64(Euclidean(start, goal)), 4))
	//fmt.Printf("presize: %v\n", presize)
	costs := make(map[Position]float32, presize)
	costs[start] = 0

	frontier := astar.queueBuilder()
	frontier.Push(start, 0)
	parentChain := make(map[Position]Position, presize)

	var iter uint
	var curr Position

	for !frontier.Empty() {
		iter++
		curr = frontier.Pop().(Position)

		if curr == goal {
			break
		}

		neighbors := world.GetNeighbors(curr, astar.Diagonal)
		for _, node := range neighbors {
			tileCost := world.GetCost(node)
			newCost := costs[curr] + tileCost
			prevCost, exists := costs[node]
			if !exists || newCost < prevCost {
				costs[node] = newCost
				parentChain[node] = curr
				frontier.Push(node, newCost+astar.Heuristic(node, goal))
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

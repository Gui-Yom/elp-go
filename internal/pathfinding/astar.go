package pathfinding

import (
	"elp-go/internal/fastmap"
	"elp-go/internal/queue"
	"elp-go/internal/world"
	"math"
	"time"
)

type Astar struct {
	Diagonal     bool
	Heuristic    Heuristic
	queueBuilder func() queue.PriorityQueue
}

var _ Pathfinder = (*Astar)(nil)

func (astar Astar) FindPath(w *world.World, start world.Position, goal world.Position) ([]world.Position, Stats) {
	startTime := time.Now()

	// Astar explores in the direction of the goal, we model it by a segment
	// We multiply that number to account for the fact there is walls
	presize := int(math.Max(2*math.Sqrt2*float64(Euclidean(start, goal)), 4))
	//fmt.Printf("presize: %v\n", presize)
	costs := fastmap.New(presize, 0.9)
	costs.PutCost(start, 0.0)

	frontier := astar.queueBuilder()
	frontier.Push(start, 0)
	parentChain := fastmap.New(presize, 0.9)

	var iter uint
	var curr world.Position

	neighbors := make([]world.Position, 8)

	for !frontier.Empty() {
		curr = frontier.Pop()

		if curr == goal {
			break
		}

		iter++

		size := w.GetNeighbors(curr, astar.Diagonal, neighbors)
		for i := 0; i < size; i++ {
			node := neighbors[i]
			tileCost := w.GetCost(node)
			currCost, _ := costs.GetCost(curr)
			newCost := currCost + tileCost
			prevCost, exists := costs.GetCost(node)
			if !exists || newCost < prevCost {
				costs.PutCost(node, newCost)
				parentChain.PutPos(node, curr)
				frontier.Push(node, newCost+astar.Heuristic(node, goal))
			}
		}
	}
	//fmt.Printf("costLen: %v, chainLen: %v\n", len(costs), len(parentChain))
	if curr != goal {
		return nil, Stats{Iterations: iter, Duration: time.Now().Sub(startTime), PresizeAccuracy: float64(costs.Size()) / float64(presize) * 100}
	}
	path := makePath(parentChain, start, goal)
	return path, Stats{Iterations: iter, Duration: time.Now().Sub(startTime), Cost: pathCost(w, path), PresizeAccuracy: float64(costs.Size()) / float64(presize) * 100}
}

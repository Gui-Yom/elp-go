package pathfinding

import (
	"elp-go/internal/fastmap"
	"elp-go/internal/queue"
	"elp-go/internal/world"
	"math"
	"time"
)

type Dijkstra struct {
	Diagonal     bool
	queueBuilder func() queue.PriorityQueue
}

// Implementation is implicit, thanks Go
var _ Pathfinder = (*Dijkstra)(nil)

func (dijk Dijkstra) FindPath(w *world.World, start world.Position, goal world.Position) ([]world.Position, Stats) {
	startTime := time.Now()

	// Dijkstra explores everywhere, we model it by a disc pi*r^2
	// We divide by 8 to account for the fact that most of the disc is out of bounds or walls
	presize := int(math.Max(math.Pi*float64(EuclideanSq(start, goal))/8.0, 4))
	//fmt.Printf("presize: %v\n", presize)
	costs := fastmap.New(presize, 0.9)
	costs.PutCost(start, 0.0)

	frontier := dijk.queueBuilder()
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

		size := w.GetNeighbors(curr, dijk.Diagonal, neighbors)
		for i := 0; i < size; i++ {
			node := neighbors[i]
			tileCost := w.GetCost(node)
			currCost, _ := costs.GetCost(curr)
			newCost := currCost + tileCost
			prevCost, exists := costs.GetCost(node)
			if !exists || newCost < prevCost {
				costs.PutCost(node, newCost)
				parentChain.PutPos(node, curr)
				frontier.Push(node, newCost)
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

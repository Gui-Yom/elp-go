package internal

import (
	"time"
)

type Dijkstra struct {
	Diagonal bool
}

// Implementation is implicit, thanks Go
var _ Pathfinder = (*Dijkstra)(nil)

func (dijk Dijkstra) FindPath(carte *Carte, start Position, goal Position) ([]Position, Stats) {
	startTime := time.Now()

	costs := make(map[Position]float32)
	costs[start] = 0

	frontier := PriorityQueue{}
	frontier.push(start, 0)
	parentChain := make(map[Position]Position)

	var iter uint
	var curr Position

	for !frontier.empty() {
		curr = frontier.pop().(Position)

		if curr == goal {
			break
		}

		iter++

		neighbors := carte.GetNeighbors(curr, dijk.Diagonal)
		for _, node := range neighbors {
			tileCost := carte.GetTile(node).Cost
			newCost := costs[curr] + tileCost
			prevCost, exists := costs[node]
			if !exists || newCost < prevCost {
				costs[node] = newCost
				parentChain[node] = curr
				frontier.push(node, newCost)
			}
		}
	}
	if curr != goal {
		return nil, Stats{Iterations: iter, Duration: time.Now().Sub(startTime)}
	}
	path := makePath(parentChain, start, goal)
	return path, Stats{Iterations: iter, Duration: time.Now().Sub(startTime), Cost: pathCost(carte, path)}
}

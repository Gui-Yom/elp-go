package internal

import (
	"time"
)

type Astar struct {
	diagonal  bool
	heuristic Heuristic
}

var _ Pathfinder = (*Astar)(nil)

func (astar Astar) FindPath(carte *Carte, start Position, goal Position) ([]Position, Stats) {
	startTime := time.Now()

	costs := make(map[Position]float32)
	costs[start] = 0

	frontier := PriorityQueue{}
	frontier.push(start, 0)
	parentChain := make(map[Position]Position)

	var iter uint
	var curr Position

	for !frontier.empty() {
		iter++
		curr = frontier.pop().(Position)

		if curr == goal {
			break
		}

		neighbors := carte.GetNeighbors(curr, astar.diagonal)
		for _, node := range neighbors {
			tileCost := carte.GetTile(node).Cost
			newCost := costs[curr] + tileCost
			prevCost, exists := costs[node]
			if !exists || newCost < prevCost {
				costs[node] = newCost
				parentChain[node] = curr
				frontier.push(node, newCost+astar.heuristic(node, goal))
			}
		}
	}
	if curr != goal {
		return nil, Stats{Iterations: iter, Duration: time.Now().Sub(startTime)}
	}
	path := makePath(parentChain, start, goal)
	return path, Stats{Iterations: iter, Duration: time.Now().Sub(startTime), Cost: pathCost(carte, path)}
}

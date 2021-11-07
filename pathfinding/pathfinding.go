package pathfinding

import (
	"elp-go/scenario"
	"fmt"
	"time"
)

type Pathfinder interface {
	path(carte *scenario.Carte, start scenario.Position, goal scenario.Position) ([]scenario.Position, Stats)
}

type Stats struct {
	// Number of iterations
	Iterations uint
	// Total duration of the pathfinding round
	Duration time.Duration
	// Total cost of the found path
	Cost float32
}

func (s Stats) String() string {
	return fmt.Sprintf("Stats{Iterations: %v, Duration: %v µs, Cost: %v}", s.Iterations, s.Duration.Microseconds(), s.Cost)
}

// makePath creates a path (position slice) from a parent chain and the 2 path ends
func makePath(parentChain map[scenario.Position]scenario.Position, start scenario.Position, goal scenario.Position) []scenario.Position {
	path := make([]scenario.Position, 1, len(parentChain))
	path[0] = goal
	// Follow the parentChain to create the path from the goal
	for curr := goal; curr != start; {
		n := parentChain[curr]
		path = append(path, n)
		curr = n
	}
	// Reverse the path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func pathCost(carte *scenario.Carte, path []scenario.Position) (cost float32) {
	for _, p := range path[1:] {
		cost += carte.GetTile(p).Cost
	}
	return cost
}

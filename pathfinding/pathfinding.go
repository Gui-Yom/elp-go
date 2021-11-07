package pathfinding

import (
	"elp-go/scenario"
	"log"
)

type Pathfinder interface {
	path(start scenario.Position, goal scenario.Position) []scenario.Position
}

// makePath creates a path (position slice) from a parent chain and the 2 path ends
func makePath(parentChain map[scenario.Position]scenario.Position, start scenario.Position, goal scenario.Position) []scenario.Position {

	log.Printf("len: %v", len(parentChain))
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

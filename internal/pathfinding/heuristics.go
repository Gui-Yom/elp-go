package pathfinding

import (
	"elp-go/internal/world"
	"math"
)

// Heuristic An heuristic function evaluates a given position and help us choose which path to explore first.
type Heuristic func(p world.Position, goal world.Position) float64

// Manhattan Adapted when moving in 4 directions.
func Manhattan(p world.Position, goal world.Position) float64 {
	return math.Abs(float64(p.X-goal.X)) + math.Abs(float64(p.Y-goal.Y))
}

// Euclidean Adapted when moving in 8 directions.
func Euclidean(p world.Position, goal world.Position) float64 {
	return math.Sqrt(EuclideanSq(p, goal))
}

// EuclideanSq Same as Euclidean but save us a root square.
func EuclideanSq(p world.Position, goal world.Position) float64 {
	dx := math.Abs(float64(p.X - goal.X))
	dy := math.Abs(float64(p.Y - goal.Y))
	return dx*dx + dy*dy
}

// Chebyshev Optimal heuristic when using an 8-neighborhood.
func Chebyshev(p world.Position, goal world.Position) float64 {
	dx := float64(p.X - goal.X)
	dy := float64(p.Y - goal.Y)
	return dx + dy - math.Min(dx, dy)
}

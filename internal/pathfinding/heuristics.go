package pathfinding

import (
	"math"
)

type Heuristic func(p Position, goal Position) float32

func Manhattan(p Position, goal Position) float32 {
	return float32(math.Abs(float64(p.X-goal.X)) + math.Abs(float64(p.Y-goal.Y)))
}

func Euclidean(p Position, goal Position) float32 {
	return float32(math.Sqrt(float64(EuclideanSq(p, goal))))
}

func EuclideanSq(p Position, goal Position) float32 {
	dx := math.Abs(float64(p.X - goal.X))
	dy := math.Abs(float64(p.Y - goal.Y))
	return float32(dx*dx + dy*dy)
}

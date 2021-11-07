package pathfinding

import (
	"elp-go/scenario"
	"math"
)

type Heuristic func(p scenario.Position, goal scenario.Position) float32

func Manhattan(p scenario.Position, goal scenario.Position) float32 {
	return float32(math.Abs(float64(p.X-goal.X)) + math.Abs(float64(p.Y-goal.Y)))
}

func Euclidean(p scenario.Position, goal scenario.Position) float32 {
	dx := math.Abs(float64(p.X - goal.X))
	dy := math.Abs(float64(p.Y - goal.Y))
	return float32(math.Sqrt(dx*dx + dy*dy))
}

package internal

import "elp-go/internal/pathfinding"

type Scenario struct {
	Carte            *pathfinding.Carte
	DiagonalMovement bool
	Tasks            []interface{}
	NumAgents        uint32
}

type CompletedTask struct {
	AgentId uint32
	Path    []pathfinding.Position
}

type ScenarioResult struct {
	Completed []CompletedTask
}

type Task interface {
	a()
}

type MoveTask struct {
	Goal pathfinding.Position
}

func (this MoveTask) a() {

}

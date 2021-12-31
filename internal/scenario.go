package internal

import (
	"elp-go/internal/world"
	"fmt"
)

type Scenario struct {
	World            *world.World
	DiagonalMovement bool
	Tasks            []interface{}
	NumAgents        uint32
}

type CompletedTask struct {
	AgentId uint32
	Path    []world.Position
}

type ScenarioResult struct {
	Completed []CompletedTask
}

func (res ScenarioResult) String() string {
	s := "ScenarioResult\n  Tasks:\n"
	for _, t := range res.Completed {
		s = fmt.Sprintf("%s    %v -> %v\n", s, t.AgentId, t.Path)
	}
	return s
}

type Task interface {
	Execute(agent *Agent)
}

type MoveTask struct {
	Goal world.Position
}

func (this MoveTask) Execute(agent *Agent) {

}

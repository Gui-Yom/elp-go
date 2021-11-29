package internal

type Scenario struct {
	Carte            *Carte
	DiagonalMovement bool
	Tasks            []interface{}
	NumAgents        uint32
}

type CompletedTask struct {
	AgentId uint32
	Path    []Position
}

type ScenarioResult struct {
	Completed []CompletedTask
}

type Task interface {
	a()
}

type MoveTask struct {
	Goal Position
}

func (this MoveTask) a() {

}

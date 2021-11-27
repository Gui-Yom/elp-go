package scenario

import "container/list"

type Scenario struct {
	Carte            *Carte
	DiagonalMovement bool
	Tasks            list.List
	NumAgents        int
}

func (s *Scenario) PopTask() Task {
	elem := s.Tasks.Front()
	if elem == nil {
		return nil
	} else {
		return s.Tasks.Remove(elem).(Task)
	}
}

type CompletedTask struct {
}

type ScenarioResult struct {
}

type Task interface {
	a()
}

type MoveTask struct {
	goal Position
}

func (this MoveTask) a() {

}

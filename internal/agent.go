package internal

import (
	"elp-go/internal/pathfinding"
	"log"
)

type Agent struct {
	Id         uint32
	Pos        pathfinding.Position
	pathfinder pathfinding.Pathfinder
}

func NewAgent(id uint32, pos pathfinding.Position, pathfinder pathfinding.Pathfinder) Agent {
	return Agent{Id: id, Pos: pos, pathfinder: pathfinder}
}

func (a Agent) ExecuteTask(carte *pathfinding.Carte, task Task) CompletedTask {
	switch t := task.(type) {
	case MoveTask:
		log.Printf("%v -> %v", a.Id, t)
		path, _ := a.pathfinder.FindPath(carte, a.Pos, t.Goal)
		// TODO(guillaume) pass stats
		return CompletedTask{AgentId: a.Id, Path: path}
	default:
		log.Fatalf("Unimplemented task : %v", t)
	}
	return CompletedTask{}
}

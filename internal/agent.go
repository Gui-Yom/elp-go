package internal

import (
	"elp-go/internal/pathfinding"
	"elp-go/internal/world"
	"log"
)

type Agent struct {
	Id         uint
	Pos        world.Position
	pathfinder pathfinding.Pathfinder
}

func NewAgent(id uint, pos world.Position, pathfinder pathfinding.Pathfinder) Agent {
	return Agent{Id: id, Pos: pos, pathfinder: pathfinder}
}

func (a *Agent) ExecuteTask(world *world.World, task Task) CompletedTask {
	switch t := task.(type) {
	case MoveTask:
		//log.Printf("%v -> %v", a.Id, t)
		path, stats := a.pathfinder.FindPath(world, a.Pos, t.Goal)
		a.Pos = t.Goal
		return CompletedTask{AgentId: a.Id, Path: path, Stats: stats}
	default:
		log.Fatalf("Unimplemented task : %v", t)
	}
	return CompletedTask{}
}

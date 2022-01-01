package internal

import (
	"bufio"
	"elp-go/internal/pathfinding"
	"elp-go/internal/world"
	"fmt"
	"log"
	"os"
)

type Scenario struct {
	World            *world.World
	DiagonalMovement bool
	Tasks            []interface{}
	NumAgents        int
}

func LoadFromFile(filename string) Scenario {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ok := scanner.Scan()
	if !ok {
		panic("expected 'map <mapname>'")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type CompletedTask struct {
	AgentId uint
	Path    []world.Position
	Stats   pathfinding.Stats
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

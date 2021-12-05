package internal

import (
	"elp-go/internal/pathfinding"
	"elp-go/internal/queue"
	"fmt"
	"math/rand"
	"testing"
)

func genTasks(num int, width, height int) []interface{} {
	tasks := make([]interface{}, num)
	for i := 0; i < len(tasks); i++ {
		tasks[i] = MoveTask{Goal: pathfinding.Pos(rand.Intn(width), rand.Intn(height))}
	}
	return tasks
}

func testRequestHandler(t *testing.T, handler RequestHandler) {
	scen := Scenario{
		Carte:            pathfinding.NewMapEmpty(100, 100),
		DiagonalMovement: true,
		Tasks:            genTasks(8, 100, 100),
		NumAgents:        4,
	}
	result := handler(&scen, pathfinding.NewDijkstra(true, queue.NewLinked))
	fmt.Printf("result : %v", result)
}

func TestHandleRequestSeq(t *testing.T) {
	testRequestHandler(t, handleRequestSeq)
}

func TestHandleRequestPar(t *testing.T) {
	testRequestHandler(t, handleRequestPar)
}

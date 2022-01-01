package internal

import (
	"elp-go/internal/pathfinding"
	"elp-go/internal/queue"
	"elp-go/internal/world"
	"fmt"
	"math/rand"
	"testing"
)

func genTasks(num int, w *world.World) []interface{} {
	tasks := make([]interface{}, num)
	for i := 0; i < len(tasks); i++ {
		pos := world.Pos(rand.Intn(w.Width), rand.Intn(w.Height))
		for ; !w.GetTile(pos).IsTraversable(); pos = world.Pos(rand.Intn(w.Width), rand.Intn(w.Height)) {
		}
		tasks[i] = MoveTask{Goal: pos}
	}
	return tasks
}

func testRequestHandler(t *testing.T, handler RequestHandler) {
	w := world.NewWorldEmpty(100, 100)
	scen := Scenario{
		World:            w,
		DiagonalMovement: true,
		Tasks:            genTasks(8, w),
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

func benchmarkRequestHandler(b *testing.B, handler RequestHandler) {
	b.StopTimer()
	b.ResetTimer()
	pathfinder := pathfinding.NewAstar(true, pathfinding.EuclideanSq, queue.NewLinked)
	w := world.NewWorldRandom(1000, 1000, 0.2, 42)
	scenario := Scenario{
		World:            w,
		DiagonalMovement: true,
		Tasks:            genTasks(256, w),
		NumAgents:        6,
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		handler(&scenario, pathfinder)
	}
}

func BenchmarkSeq(b *testing.B) {
	benchmarkRequestHandler(b, handleRequestSeq)
}

func BenchmarkPar(b *testing.B) {
	benchmarkRequestHandler(b, handleRequestPar)
}

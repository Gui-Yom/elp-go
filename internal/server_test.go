package internal

import (
	"elp-go/internal/pathfinding"
	"elp-go/internal/queue"
	"elp-go/internal/world"
	"fmt"
	"testing"
)

func testRequestHandler(t *testing.T, handler RequestHandler) {
	w := world.NewWorldEmpty(100, 100)
	scen := Scenario{
		World:            w,
		DiagonalMovement: true,
		Tasks:            genTasks(8, w),
		Agents:           make([]world.Position, 4),
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

func benchmarkRequestHandler(b *testing.B, handler RequestHandler, numAgents int) {
	b.StopTimer()
	b.ResetTimer()
	pathfinder := pathfinding.NewAstar(true, pathfinding.EuclideanSq, queue.NewLinked)
	w := world.NewWorldRandom(1000, 1000, 0.2, 42)
	numTasks := 256.0
	scenario := Scenario{
		World:            w,
		DiagonalMovement: true,
		Tasks:            genTasks(int(numTasks), w),
		Agents:           make([]world.Position, numAgents),
	}
	var accDuration float64
	var accIterations float64
	var accCost float64
	var presizeAcc float64
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		result := handler(&scenario, pathfinder)
		b.StopTimer()
		for _, t := range result.Completed {
			accDuration += float64(t.Stats.Duration.Microseconds())
			accIterations += float64(t.Stats.Iterations)
			accCost += t.Stats.Cost
			presizeAcc += t.Stats.PresizeAccuracy
		}
	}
	b.StopTimer()
	b.ReportMetric(accDuration/float64(b.N)/numTasks, "µs/op")
	b.ReportMetric(accIterations/float64(b.N)/numTasks, "iter")
	b.ReportMetric(accCost/float64(b.N)/numTasks, "cost")
	b.ReportMetric(presizeAcc/float64(b.N)/numTasks, "presizeAccuracy%")
}

func BenchmarkSeq(b *testing.B) {
	benchmarkRequestHandler(b, handleRequestSeq, 4)
}

/*
func BenchmarkPar1(b *testing.B) {
	benchmarkRequestHandler(b, handleRequestPar, 1)
}*/

func BenchmarkPar2(b *testing.B) {
	benchmarkRequestHandler(b, handleRequestPar, 2)
}

func BenchmarkPar4(b *testing.B) {
	benchmarkRequestHandler(b, handleRequestPar, 4)
}

func BenchmarkPar6(b *testing.B) {
	benchmarkRequestHandler(b, handleRequestPar, 6)
}

func BenchmarkPar8(b *testing.B) {
	benchmarkRequestHandler(b, handleRequestPar, 8)
}

package pathfinding

import (
	"elp-go/internal/queue"
	"elp-go/internal/world"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

// TESTS //

func testPathfinderSimple(t *testing.T, pf Pathfinder) {
	carte := world.NewWorldFromString(`4x4
    
xx  
    
    `)
	log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, world.Position{X: 0}, world.Position{X: 3})
	assert.NotNil(t, path, "A path should exist")
}

func testPathfinderMapFile(t *testing.T, pf Pathfinder) {
	carte := world.NewWorldFromFile("map0.map")
	log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, world.Position{}, world.Position{X: 9, Y: 9})
	assert.NotNil(t, path, "A path should exist")
}

func testPathfinderBig(t *testing.T, pf Pathfinder) {
	carte := world.NewWorldRandom(100, 100, 0.30, 42)
	log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, world.Position{}, world.Position{X: 98, Y: 97})
	assert.NotNil(t, path, "A path should exist")
}

func testPathfinderBigger(t *testing.T, pf Pathfinder) {
	carte := world.NewWorldRandom(300, 300, 0.30, 42)
	//log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, world.Position{}, world.Position{X: 298, Y: 298})
	assert.NotNil(t, path, "A path should exist")
}

func testPathfinderBiggger(t *testing.T, pf Pathfinder) {
	carte := world.NewWorldRandom(500, 500, 0.30, 42)
	//log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, world.Position{}, world.Position{X: 498, Y: 498})
	assert.NotNil(t, path, "A path should exist")
}

func testPathfinderBiggest(t *testing.T, pf Pathfinder) {
	carte := world.NewWorldRandom(1000, 1000, 0.30, 42)
	//log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, world.Position{}, world.Position{X: 998, Y: 998})
	assert.NotNil(t, path, "A path should exist")
}

func TestDijkstraSimple(t *testing.T) {
	testPathfinderSimple(t, NewDijkstra(false, queue.NewLinked))
}

func TestDijkstraMapFile(t *testing.T) {
	testPathfinderMapFile(t, NewDijkstra(false, queue.NewLinked))
}

func TestDijkstraBig(t *testing.T) {
	testPathfinderBig(t, NewDijkstra(true, queue.NewLinked))
}

func TestDijkstraBigger(t *testing.T) {
	testPathfinderBigger(t, NewDijkstra(true, queue.NewLinked))
}

func TestDijkstraBiggger(t *testing.T) {
	testPathfinderBiggger(t, NewDijkstra(true, queue.NewLinked))
}

func TestDijkstraBiggest(t *testing.T) {
	testPathfinderBiggest(t, NewDijkstra(true, queue.NewLinked))
}

func TestAstarSimple(t *testing.T) {
	testPathfinderSimple(t, NewAstar(false, Manhattan, queue.NewLinked))
}

func TestAstarMapFile(t *testing.T) {
	testPathfinderMapFile(t, NewAstar(false, Manhattan, queue.NewLinked))
}

func TestAstarBig(t *testing.T) {
	testPathfinderBig(t, NewAstar(true, Euclidean, queue.NewLinked))
}

func TestAstarBigger(t *testing.T) {
	testPathfinderBigger(t, NewAstar(true, Euclidean, queue.NewLinked))
}

func TestAstarBiggger(t *testing.T) {
	testPathfinderBiggger(t, NewAstar(true, Euclidean, queue.NewLinked))
}

func TestAstarBiggest(t *testing.T) {
	testPathfinderBiggest(t, NewAstar(true, Euclidean, queue.NewLinked))
}

// BENCHMARKS //

func benchmarkPathfinder(pathfinder Pathfinder, carte *world.World, goal world.Position, b *testing.B) {
	b.StopTimer()
	b.ResetTimer()
	var accDuration float64
	var accIterations float64
	var accCost float64
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, stats := pathfinder.FindPath(carte, world.Position{}, goal)
		accDuration += float64(stats.Duration.Microseconds())
		accIterations += float64(stats.Iterations)
		accCost += float64(stats.Cost)
	}
	b.StopTimer()
	b.ReportMetric(accDuration/float64(b.N), "µs/op")
	b.ReportMetric(accIterations/float64(b.N), "iter/op")
	b.ReportMetric(accCost/float64(b.N), "cost/op")
}

func BenchmarkDijkstraLinked100(b *testing.B) {
	benchmarkPathfinder(
		NewDijkstra(true, queue.NewLinked),
		world.NewWorldRandom(100, 100, 0.3, 42),
		world.Pos(97, 98),
		b)
}

func BenchmarkDijkstraLinked1000(b *testing.B) {
	benchmarkPathfinder(
		NewDijkstra(true, queue.NewLinked),
		world.NewWorldRandom(1000, 1000, 0.2, 42),
		world.Pos(998, 998),
		b)
}

func BenchmarkDijkstraPairing100(b *testing.B) {
	benchmarkPathfinder(
		NewDijkstra(true, queue.NewPairing),
		world.NewWorldRandom(100, 100, 0.3, 42),
		world.Pos(98, 97),
		b)
}

func BenchmarkDijkstraPairing300(b *testing.B) {
	benchmarkPathfinder(
		NewDijkstra(true, queue.NewPairing),
		world.NewWorldRandom(300, 300, 0.2, 42),
		world.Pos(298, 298),
		b)
}

func BenchmarkAstarLinked100(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewLinked),
		world.NewWorldRandom(100, 100, 0.3, 42),
		world.Pos(98, 97),
		b)
}

func BenchmarkAstarLinked1000(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewLinked),
		world.NewWorldRandom(1000, 1000, 0.2, 42),
		world.Pos(998, 998),
		b)
}

func BenchmarkAstarLinked10000(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewLinked),
		world.NewWorldRandom(10000, 10000, 0.2, 42),
		world.Pos(9998, 9998),
		b)
}

func BenchmarkAstarPairing100(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewPairing),
		world.NewWorldRandom(100, 100, 0.3, 42),
		world.Pos(98, 97),
		b)
}

func BenchmarkAstarPairing300(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewPairing),
		world.NewWorldRandom(300, 300, 0.2, 42),
		world.Pos(298, 298),
		b)
}

func BenchmarkAstarPairing1000(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewPairing),
		world.NewWorldRandom(1000, 1000, 0.2, 42),
		world.Pos(998, 998),
		b)
}

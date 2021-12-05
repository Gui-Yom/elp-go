package pathfinding

import (
	"elp-go/internal/queue"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func testPathfinderSimple(t *testing.T, pf Pathfinder) {
	carte := NewMapFromString(`4x4
    
xx  
    
    `)
	log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, Position{X: 0}, Position{X: 3})
	assert.NotNil(t, path, "A path should exist")
}

func testPathfinderMapFile(t *testing.T, pf Pathfinder) {
	carte := NewMapFromFile("map0.map")
	log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, Position{}, Position{X: 9, Y: 9})
	assert.NotNil(t, path, "A path should exist")
}

func testPathfinderBig(t *testing.T, pf Pathfinder) {
	carte := NewMapRandom(100, 100, 0.30, 42)
	log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, Position{}, Position{X: 98, Y: 97})
	assert.NotNil(t, path, "A path should exist")
}

func testPathfinderBigger(t *testing.T, pf Pathfinder) {
	carte := NewMapRandom(300, 300, 0.30, 42)
	//log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, Position{}, Position{X: 298, Y: 298})
	assert.NotNil(t, path, "A path should exist")
}

func testPathfinderBiggger(t *testing.T, pf Pathfinder) {
	carte := NewMapRandom(500, 500, 0.30, 42)
	//log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, Position{}, Position{X: 498, Y: 498})
	assert.NotNil(t, path, "A path should exist")
}

func testPathfinderBiggest(t *testing.T, pf Pathfinder) {
	carte := NewMapRandom(1000, 1000, 0.30, 42)
	//log.Printf("map: %v", carte)
	path, _ := pf.FindPath(carte, Position{}, Position{X: 998, Y: 998})
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

func benchmarkPathfinder(pathfinder Pathfinder, carte *Carte, goal Position, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path, stats := pathfinder.FindPath(carte, Position{}, goal)
		if path == nil {
			b.Fatal("Path not found")
		} else {
			b.ReportMetric(float64(stats.Duration.Microseconds())/float64(b.N), "µs/op")
		}
	}
}

func BenchmarkDijkstraLinked100(b *testing.B) {
	benchmarkPathfinder(
		NewDijkstra(true, queue.NewLinked),
		NewMapRandom(100, 100, 0.3, 42),
		Pos(97, 98),
		b)
}

func BenchmarkDijkstraLinked300(b *testing.B) {
	benchmarkPathfinder(
		NewDijkstra(true, queue.NewLinked),
		NewMapRandom(300, 300, 0.2, 42),
		Pos(298, 298),
		b)
}

func BenchmarkDijkstraLinked500(b *testing.B) {
	benchmarkPathfinder(
		NewDijkstra(true, queue.NewLinked),
		NewMapRandom(500, 500, 0.2, 42),
		Pos(498, 498),
		b)
}

func BenchmarkDijkstraPairing100(b *testing.B) {
	benchmarkPathfinder(
		NewDijkstra(true, queue.NewPairing),
		NewMapRandom(100, 100, 0.3, 42),
		Pos(98, 97),
		b)
}

func BenchmarkDijkstraPairing300(b *testing.B) {
	benchmarkPathfinder(
		NewDijkstra(true, queue.NewPairing),
		NewMapRandom(300, 300, 0.2, 42),
		Pos(298, 298),
		b)
}

func BenchmarkDijkstraPairing500(b *testing.B) {
	benchmarkPathfinder(
		NewDijkstra(true, queue.NewPairing),
		NewMapRandom(500, 500, 0.2, 42),
		Pos(498, 498),
		b)
}

func BenchmarkAstarLinked100(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewLinked),
		NewMapRandom(100, 100, 0.3, 42),
		Pos(98, 97),
		b)
}

func BenchmarkAstarLinked300(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewLinked),
		NewMapRandom(300, 300, 0.2, 42),
		Pos(298, 298),
		b)
}

func BenchmarkAstarLinked500(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewLinked),
		NewMapRandom(500, 500, 0.2, 42),
		Pos(498, 498),
		b)
}

func BenchmarkAstarPairing100(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewPairing),
		NewMapRandom(100, 100, 0.3, 42),
		Pos(98, 97),
		b)
}

func BenchmarkAstarPairing300(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewPairing),
		NewMapRandom(300, 300, 0.2, 42),
		Pos(298, 298),
		b)
}

func BenchmarkAstarPairing500(b *testing.B) {
	benchmarkPathfinder(
		NewAstar(true, EuclideanSq, queue.NewPairing),
		NewMapRandom(500, 500, 0.2, 42),
		Pos(498, 498),
		b)
}

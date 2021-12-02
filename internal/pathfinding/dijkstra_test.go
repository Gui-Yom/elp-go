package pathfinding

import (
	"log"
	"testing"
)

func TestDijkstra(t *testing.T) {
	dijk := Dijkstra{}
	carte := NewMapFromString(`4x4
    
xx  
    
    `)
	log.Printf("map: %v", carte)
	path, _ := dijk.FindPath(carte, Position{X: 0}, Position{X: 3})
	log.Printf("path: %v", path)
}

func TestDijkstraMap0(t *testing.T) {
	dijk := Dijkstra{}
	carte := NewMapFromFile("../map0.map")
	log.Printf("map: %v", carte)
	path, stats := dijk.FindPath(carte, Position{}, Position{X: 9, Y: 9})
	log.Printf("path: %v", path)
	log.Printf("stats: %v", stats)
}

func TestDijkstraBigMap(t *testing.T) {
	dijk := Dijkstra{Diagonal: true}
	carte := NewMapRandom(100, 100, 0.30, 42)
	log.Printf("map: %v", carte)
	path, stats := dijk.FindPath(carte, Position{}, Position{X: 98, Y: 97})
	log.Printf("path: %v", path)
	log.Printf("stats: %v", stats)
}

func TestDijkstraBiggerMap(t *testing.T) {
	dijk := Dijkstra{Diagonal: true}
	carte := NewMapRandom(300, 300, 0.30, 42)
	log.Printf("map: %v", carte)
	path, stats := dijk.FindPath(carte, Position{}, Position{X: 298, Y: 298})
	log.Printf("path: %v", path)
	log.Printf("stats: %v", stats)
}

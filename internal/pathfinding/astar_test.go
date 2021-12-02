package pathfinding

import (
	"log"
	"testing"
)

func TestAstar(t *testing.T) {
	astar := Astar{heuristic: Manhattan}
	carte := NewMapFromString(`4x4
    
xx  
    
    `)
	log.Printf("map: %v", carte)
	path, _ := astar.FindPath(carte, Position{X: 0}, Position{X: 3})
	log.Printf("path: %v", path)
}

func TestAstarMap0(t *testing.T) {
	astar := Astar{heuristic: Manhattan}
	carte := NewMapFromFile("../map0.map")
	log.Printf("map: %v", carte)
	path, stats := astar.FindPath(carte, Position{}, Position{X: 9, Y: 9})
	log.Printf("path: %v", path)
	log.Printf("stats: %v", stats)
}

func TestAstarBigMap(t *testing.T) {
	astar := Astar{diagonal: true, heuristic: Euclidean}
	carte := NewMapRandom(100, 100, 0.30, 42)
	log.Printf("map: %v", carte)
	path, stats := astar.FindPath(carte, Position{}, Position{X: 98, Y: 97})
	log.Printf("path: %v", path)
	log.Printf("stats: %v", stats)
}

func TestAstarBiggerMap(t *testing.T) {
	astar := Astar{diagonal: true, heuristic: Euclidean}
	carte := NewMapRandom(300, 300, 0.30, 42)
	log.Printf("map: %v", carte)
	path, stats := astar.FindPath(carte, Position{}, Position{X: 298, Y: 298})
	log.Printf("path: %v", path)
	log.Printf("stats: %v", stats)
}

package pathfinding

import (
	"elp-go/scenario"
	"log"
	"testing"
)

func TestDijkstra(t *testing.T) {
	dijk := Dijkstra{}
	carte := scenario.ReadMapFromString(`4x4
    
xx  
    
    `)
	log.Printf("map:\n%v", carte)
	path, _ := dijk.path(carte, scenario.Position{X: 0}, scenario.Position{X: 3})
	log.Printf("path: %v", path)
}

func TestDijkstraMap0(t *testing.T) {
	dijk := Dijkstra{}
	carte := scenario.ReadMapFromFile("../map0.map")
	log.Printf("map:\n%v", carte)
	path, stats := dijk.path(carte, scenario.Position{}, scenario.Position{X: 9, Y: 9})
	log.Printf("path: %v", path)
	log.Printf("stats: %v", stats)
}

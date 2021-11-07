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
	path := dijk.path(carte, scenario.Position{X: 0}, scenario.Position{X: 3})
	log.Printf("path: %v", path)
}

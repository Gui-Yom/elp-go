package pathfinding

import (
	"elp-go/scenario"
	"log"
	"testing"
)

func TestAstar(t *testing.T) {
	astar := Astar{heuristic: Manhattan}
	carte := scenario.ReadMapFromString(`4x4
    
xx  
    
    `)
	log.Printf("map:\n%v", carte)
	path := astar.path(carte, scenario.Position{X: 0}, scenario.Position{X: 3})
	log.Printf("path: %v", path)
}

package scenario

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Tile struct {
	id   uint8
	cost int
}

// No enums, thank you Go
var (
	TILE_EMPTY = Tile{id: ' ', cost: 1}
	TILE_WALL  = Tile{id: 'x', cost: -1}
	TILE_GOAL  = Tile{id: 'G', cost: 1}
	// Apparently Go is stupid and just return a copy of the value each time we access it instead of a reference, thank you Go.
	TILES = map[uint8]Tile{TILE_EMPTY.id: TILE_EMPTY, TILE_WALL.id: TILE_WALL, TILE_GOAL.id: TILE_GOAL}
)

type Carte struct {
	Row   int
	Col   int
	inner []uint8
}

func (c *Carte) GetRaw(x int, y int) uint8 {
	return c.inner[x*c.Row+y]
}

func (c *Carte) GetTile(x int, y int) Tile {
	return TILES[c.GetRaw(x, y)]
}

func ReadMap(name string) Carte {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	// Schedule file close immediately
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		log.Fatal("map file is invalid")
	}
	// Parse map size (<row>x<col>)
	bits := strings.Split(scanner.Text(), "x")
	row, _ := strconv.Atoi(bits[0])
	col, _ := strconv.Atoi(bits[1])

	tab := make([]uint8, row*col, row*col)
	for i := 0; scanner.Scan(); i++ {
		for j, char := range strings.TrimRight(scanner.Text(), " \t\n\r") {

			// We only take the lowest part of the unicode codepoint for ascii
			id := uint8(char & 0x000000FF)
			//log.Printf("char: %v, id: %v", char, id)

			// Is end of map line ?
			if id != '|' {
				// Is the char is a valid identifier ?
				if _, exists := TILES[id]; !exists {
					log.Fatal("Invalid tile : '", id, "'")
				}
				tab[i*row+j] = id
			}
		}
	}

	// Check for an error
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return Carte{Row: row, Col: col, inner: tab}
}

func (c Carte) String() string {
	var s = ""
	for i := 0; i < c.Row; i++ {
		for j := 0; j < c.Col; j++ {
			s += string(rune(c.GetRaw(i, j)))
		}
		s += "\n"
	}
	return s
}

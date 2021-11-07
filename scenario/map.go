package scenario

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	X int
	Y int
}

func (p Position) Plus(o Position) Position {
	return Position{X: p.X + o.X, Y: p.Y + o.Y}
}

func (p Position) String() string {
	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}

type Tile struct {
	Id   uint8
	Cost float32
}

// No enums, thank you Go
var (
	TILE_EMPTY         = &Tile{Id: ' ', Cost: 1}
	TILE_WALL          = &Tile{Id: 'x', Cost: -1}
	TILE_GOAL          = &Tile{Id: 'G', Cost: 1}
	TILE_CONVEYOR_BELT = &Tile{Id: 'C', Cost: 0.5}
	TILE_STAIRS        = &Tile{Id: '[', Cost: 2}
	TILE_LADDER        = &Tile{Id: '#', Cost: 2}
	TILE_SLIPPERY_ROCK = &Tile{Id: '(', Cost: 3}
	TILE_HILL          = &Tile{Id: 'H', Cost: 3}
	TILE_STREAM        = &Tile{Id: '~', Cost: 4}
	TILE_HOLE          = &Tile{Id: 'o', Cost: 4}
	TILE_CAVE          = &Tile{Id: 'C', Cost: 4}

	TILES = map[uint8]*Tile{
		TILE_EMPTY.Id:         TILE_EMPTY,
		TILE_WALL.Id:          TILE_WALL,
		TILE_GOAL.Id:          TILE_GOAL,
		TILE_CONVEYOR_BELT.Id: TILE_CONVEYOR_BELT,
	}
)

// IsTraversable a tile is traversable if its cost is > 0
func (t *Tile) IsTraversable() bool {
	return t.Cost > 0
}

type Carte struct {
	Row   int
	Col   int
	Inner []uint8
}

func (c *Carte) IsInBounds(p Position) bool {
	return p.X >= 0 && p.X < c.Col && p.Y >= 0 && p.Y < c.Row
}

func (c *Carte) GetRaw(x int, y int) uint8 {
	return c.Inner[x*c.Row+y]
}

func (c *Carte) GetTile(p Position) *Tile {
	return TILES[c.GetRaw(p.X, p.Y)]
}

// GetNeighbors returns traversable tiles around position (x, y)
func (c *Carte) GetNeighbors(p Position, diagonal bool) []Position {
	// No functional programming, thanks Go
	var arr []Position
	var offsets []Position
	if diagonal {
		arr = make([]Position, 0, 8)
		offsets = []Position{
			{X: -1, Y: 1}, {Y: 1}, {X: 1, Y: 1},
			{X: -1}, {X: 1},
			{X: -1, Y: -1}, {Y: -1}, {X: 1, Y: -1}}
	} else {
		arr = make([]Position, 0, 4)
		offsets = []Position{{X: 1}, {Y: 1}, {X: -1}, {Y: -1}}
	}
	for _, offset := range offsets {
		n := p.Plus(offset)
		if c.IsInBounds(n) && c.GetTile(n).IsTraversable() {
			arr = append(arr, n)
		}
	}
	return arr
}

func ReadMapFromFile(name string) *Carte {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	// Schedule file close immediately
	defer file.Close()

	return ReadMap(file)
}

func ReadMapFromString(mapText string) *Carte {
	return ReadMap(strings.NewReader(mapText))
}

func ReadMap(r io.Reader) *Carte {
	scanner := bufio.NewScanner(r)

	if !scanner.Scan() {
		log.Fatal("map file is invalid")
	}
	// Parse map size (<row>x<col>)
	bits := strings.Split(scanner.Text(), "x")
	row, _ := strconv.Atoi(bits[0])
	col, _ := strconv.Atoi(bits[1])

	tab := make([]uint8, row*col, row*col)
	for i := 0; scanner.Scan(); i++ {
		for j, char := range strings.TrimRight(scanner.Text(), "\t\n\r") {

			// We only take the lowest part of the unicode codepoint for ascii
			id := uint8(char & 0x000000FF)

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

	return &Carte{Row: row, Col: col, Inner: tab}
}

func RandomMap(row int, col int, fill float32, seed int64) *Carte {
	rand := rand.New(rand.NewSource(seed))
	inner := make([]uint8, row*col)
	for i := 0; i < row*col; i++ {
		if rand.Float32() < fill {
			inner[i] = 'x'
		} else {
			inner[i] = ' '
		}
	}
	return &Carte{Row: row, Col: col, Inner: inner}
}

func (c Carte) String() string {
	var s = fmt.Sprintf("%vx%v\n", c.Row, c.Col)
	for i := 0; i < c.Row; i++ {
		for j := 0; j < c.Col; j++ {
			s += string(rune(c.GetRaw(i, j)))
		}
		s += "\n"
	}
	return s
}

package world

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

// Position A position in the discrete 2D world
type Position struct {
	X int
	Y int
}

func Pos(x, y int) Position {
	return Position{X: x, Y: y}
}

func (p Position) Plus(o Position) Position {
	return Position{X: p.X + o.X, Y: p.Y + o.Y}
}

func (p Position) String() string {
	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}

// Tile A tile describes the conditions of a terrain (e.g. cost to traverse).
// A negative cost indicates this tile is not traversable.
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

// World The discrete 2D world our agents operate in.
type World struct {
	Width  int
	Height int
	Inner  []uint8
}

func (w *World) IsInBounds(p Position) bool {
	return p.X >= 0 && p.X < w.Width && p.Y >= 0 && p.Y < w.Height
}

func (w *World) GetRaw(x int, y int) uint8 {
	return w.Inner[y*w.Width+x]
}

func (w *World) GetTile(p Position) *Tile {
	return TILES[w.GetRaw(p.X, p.Y)]
}

func (w *World) GetCost(p Position) (cost float32) {
	// Switch expressions ?
	switch w.GetRaw(p.X, p.Y) {
	case TILE_EMPTY.Id:
		cost = TILE_EMPTY.Cost
	case TILE_WALL.Id:
		cost = TILE_WALL.Cost
	case TILE_GOAL.Id:
		cost = TILE_GOAL.Cost
	case TILE_CONVEYOR_BELT.Id:
		cost = TILE_CONVEYOR_BELT.Cost
	default:
		panic("Unknown tile")
	}
	return cost
}

// GetNeighbors returns traversable tiles around position (x, y)
func (w *World) GetNeighbors(p Position, diagonal bool) (pos []Position) {
	// No functional programming, thanks Go
	// I wonder if there is a way to not allocate, this function gets called for each tile of our world.
	// This may probably be calculated ahead of time since the world won't change.
	var offsets []Position
	if diagonal {
		pos = make([]Position, 0, 8)
		offsets = []Position{
			{X: -1, Y: 1}, {Y: 1}, {X: 1, Y: 1},
			{X: -1}, {X: 1},
			{X: -1, Y: -1}, {Y: -1}, {X: 1, Y: -1}}
	} else {
		pos = make([]Position, 0, 4)
		offsets = []Position{{X: 1}, {Y: 1}, {X: -1}, {Y: -1}}
	}
	for _, offset := range offsets {
		n := p.Plus(offset)
		if w.IsInBounds(n) && w.GetCost(n) > 0 {
			pos = append(pos, n)
		}
	}
	return pos
}

func NewMapFromFile(name string) *World {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	// Schedule file close immediately
	defer file.Close()

	return NewMap(file)
}

func NewMapFromString(mapText string) *World {
	return NewMap(strings.NewReader(mapText))
}

func NewMap(r io.Reader) *World {
	scanner := bufio.NewScanner(r)

	if !scanner.Scan() {
		log.Fatal("map file is invalid")
	}
	// Parse map size (<width>x<height>)
	bits := strings.Split(scanner.Text(), "x")
	width, _ := strconv.Atoi(bits[0])
	height, _ := strconv.Atoi(bits[1])

	tab := make([]uint8, width*height)
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
				tab[i*width+j] = id
			}
		}
	}

	// Check for an error
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return &World{Width: width, Height: height, Inner: tab}
}

func NewMapRandom(width int, height int, fill float32, seed int64) *World {
	rand := rand.New(rand.NewSource(seed))
	inner := make([]uint8, width*height)
	for i := 0; i < width*height; i++ {
		if rand.Float32() < fill {
			inner[i] = 'x'
		} else {
			inner[i] = ' '
		}
	}
	return &World{Width: width, Height: height, Inner: inner}
}

func NewMapEmpty(width, height int) *World {
	inner := make([]uint8, width*height)
	for i := range inner {
		inner[i] = ' '
	}
	return &World{Width: width, Height: height, Inner: inner}
}

func (w World) String() string {
	var s = fmt.Sprintf("%vx%v\n", w.Width, w.Height)
	for j := 0; j < w.Height; j++ {
		for i := 0; i < w.Width; i++ {
			s += string(rune(w.GetRaw(i, j)))
		}
		s += "\n"
	}
	return s
}

func (w World) SaveToFile(name string) {
	if file, err := os.Create(name); err != nil {
		log.Fatal(err)
	} else {
		defer file.Close()
		if _, err := file.WriteString(w.String()); err != nil {
			log.Fatal(err)
		}
	}
}

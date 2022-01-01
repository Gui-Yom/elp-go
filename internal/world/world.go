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
	X int32
	Y int32
}

func Pos(x, y int) Position {
	return Position{X: int32(x), Y: int32(y)}
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
	Cost float64
}

// No enums, thank you Go
var (
	TILE_EMPTY         = &Tile{Id: ' ', Cost: 1}
	TILE_WALL          = &Tile{Id: 'x', Cost: -1}
	TILE_CONVEYOR_BELT = &Tile{Id: 'c', Cost: 0.5}
	TILE_SAND          = &Tile{Id: 's', Cost: 2}

	TILES = map[uint8]*Tile{
		TILE_EMPTY.Id:         TILE_EMPTY,
		TILE_WALL.Id:          TILE_WALL,
		TILE_CONVEYOR_BELT.Id: TILE_CONVEYOR_BELT,
		TILE_SAND.Id:          TILE_SAND,
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
	return p.X >= 0 && p.X < int32(w.Width) && p.Y >= 0 && p.Y < int32(w.Height)
}

func (w *World) GetRaw(x int, y int) uint8 {
	return w.Inner[y*w.Width+x]
}

func (w *World) GetTile(p Position) *Tile {
	return TILES[w.GetRaw(int(p.X), int(p.Y))]
}

func (w *World) GetCost(p Position) (cost float64) {
	// We don't use the TILES map because we get a nice speedup by using a switch
	// Switch expressions ?
	switch w.GetRaw(int(p.X), int(p.Y)) {
	case TILE_EMPTY.Id:
		cost = TILE_EMPTY.Cost
	case TILE_WALL.Id:
		cost = TILE_WALL.Cost
	case TILE_CONVEYOR_BELT.Id:
		cost = TILE_CONVEYOR_BELT.Cost
	case TILE_SAND.Id:
		cost = TILE_SAND.Cost
	default:
		panic("Unknown tile")
	}
	return cost
}

// This should be const but it isn't
var offsets8 = []Position{
	{X: -1, Y: 1}, {Y: 1}, {X: 1, Y: 1},
	{X: -1}, {X: 1},
	{X: -1, Y: -1}, {Y: -1}, {X: 1, Y: -1}}

var offsets4 = []Position{{X: 1}, {Y: 1}, {X: -1}, {Y: -1}}

// GetNeighbors returns traversable tiles around position (x, y)
// neighbors is an out param (with the right size), so we don't allocate repeatedly
// This should probably be expanded ahead of time since the world won't change.
func (w *World) GetNeighbors(p Position, diagonal bool, neighbors []Position) int {
	// If expressions ?
	var offsets []Position
	if diagonal {
		offsets = offsets8
	} else {
		offsets = offsets4
	}
	i := 0
	for _, offset := range offsets {
		n := p.Plus(offset)
		if w.IsInBounds(n) && w.GetCost(n) > 0 {
			neighbors[i] = n
			i++
		}
	}
	return i
}

func NewWorldFromFile(name string) *World {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	// That is actually a nice pattern
	defer file.Close()

	return NewWorld(file)
}

func NewWorldFromString(mapText string) *World {
	return NewWorld(strings.NewReader(mapText))
}

func NewWorld(r io.Reader) *World {
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
			// The character is directly the tile id
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

// NewWorldRandom
// fill increases the quantity of walls
func NewWorldRandom(width int, height int, fill float32, seed int64) *World {
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

func NewWorldEmpty(width, height int) *World {
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

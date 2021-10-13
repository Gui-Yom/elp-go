package scenario

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Carte struct {
	Row   int
	Col   int
	inner []uint8
}

func (c *Carte) get_raw(x int, y int) uint8 {
	return c.inner[x*c.Row+y]
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
		for j, char := range scanner.Text() {
			switch char {
			case ' ':
				tab[i*row+j] = 1
			case 'x':
				tab[i*row+j] = 9
			case 'G':
				tab[i*row+j] = 0
			case '|':
				// Protection for trailing whitespaces
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
	// TODO(guillaume) use tiles
	var s = ""
	for i := 0; i < c.Row; i++ {
		for j := 0; j < c.Col; j++ {
			s += fmt.Sprint(c.get_raw(i, j))
		}
		s += "\n"
	}
	return s
}

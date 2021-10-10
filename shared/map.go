package shared

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Carte struct {
	inner [][]int
}

func (c Carte) Size() (int, int) {
	return len(c.inner), len(c.inner[0])
}

func ReadMap(name string) Carte {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	// Schedule file close immediately
	defer file.Close()

	scanner := bufio.NewScanner(file)

	tab := make([][]int, 0, 4)
	for scanner.Scan() {
		inner := make([]int, 0, 4)
		for _, char := range scanner.Text() {
			valeur, _ := strconv.ParseInt(string(char), 10, 32)
			inner = append(inner, int(valeur))
		}
		tab = append(tab, inner)
	}

	// Check for an error
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, line := range tab {
		fmt.Printf("%v\n", line)
	}

	return Carte{inner: tab}
}

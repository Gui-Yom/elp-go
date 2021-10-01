package main

import "fmt"

func main() {
    carte := read_map("map.txt")
    ligne, colonne := carte.size()
    fmt.Printf("taille : %v, %v\n", ligne, colonne)
}

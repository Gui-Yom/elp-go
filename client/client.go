package client

import (
	"bufio"
	_map "elp-go/map"
	"fmt"
	"log"
	"net"
)

// Start Main func when running a client
func Start(remote string, port int) {

	client, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.ParseIP(remote), Port: port})
	if err != nil {
		log.Fatal(err)
	}
	data, err := bufio.NewReader(client).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("received : %v\n", data)

	carte := _map.ReadMap("map.txt")
	ligne, colonne := carte.Size()
	fmt.Printf("taille : %v, %v\n", ligne, colonne)
}

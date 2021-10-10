package client

import (
	_map "elp-go/map"
	"elp-go/shared"
	"fmt"
	"log"
	"net"
)

// Start Main func when running a client
func Start(addr string, port int) {

	client, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.ParseIP(addr), Port: port})
	if err != nil {
		log.Fatal(err)
	}
	remote := shared.NewRemote(client)
	var str string
	remote.Recv(&str)
	fmt.Printf("received : %v", str)

	carte := _map.ReadMap("map.txt")
	ligne, colonne := carte.Size()
	fmt.Printf("taille : %v, %v\n", ligne, colonne)
}

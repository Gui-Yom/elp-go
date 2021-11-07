package main

import (
	"elp-go/scenario"
	"fmt"
	"log"
	"net"
)

// StartClient Main func when running a client
func StartClient(addr string, port int) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.ParseIP(addr), Port: port})
	if err != nil {
		log.Fatal(err)
	}

	client := NewRemote(conn)
	defer client.Close()

	carte := scenario.ReadMapFromFile("map0.map")
	fmt.Printf("%v\n", carte)
	scenario := scenario.Scenario{Carte: *carte}

	err = client.Send(&scenario)
	if err != nil {
		log.Fatal(err)
	}
}

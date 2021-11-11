package main

import (
	"elp-go/scenario"
	"log"
	"net"
)

// StartServer Main func when running a server
func StartServer(port int) {
	server, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		client := NewRemote(conn)
		var scenario scenario.Scenario
		err = client.Recv(&scenario)
		if err != nil {
			log.Fatal(err)
		}
		handleRequest(&scenario)
	}
}

func handleRequest(scenario *scenario.Scenario) { //la fonction qui est appelée quand on reçoit une requête du client
	log.Printf("scenario: %v", scenario.Carte.Inner)
	log.Printf("scenario: %v", scenario.Carte)
}

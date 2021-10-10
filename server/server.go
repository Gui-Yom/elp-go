package server

import (
	"elp-go/shared"
	"fmt"
	"log"
	"net"
)

// Start Main func when running a server
func Start(port int) {
	server, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server is ready (%v)\n", server.Addr().String())
	// The server just wait for connections
	for {
		client, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			remote := shared.NewRemote(client)
			var msg shared.Message
			err := remote.Recv(&msg)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v\n", msg)
		}()
	}
}

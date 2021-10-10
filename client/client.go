package client

import (
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
	err = remote.Send(&shared.Message{Msg: "Hello there !"})
	remote.Close()
	if err != nil {
		log.Fatal(err)
	}

	carte := shared.ReadMap("map0.map")
	fmt.Printf("%v\n", carte)
}

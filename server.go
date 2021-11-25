package main

import (
	"elp-go/scenario"
	"log"
	"net"
	"sync"
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

func handleRequest(scen *scenario.Scenario) {
	sequential := false

	if sequential {
		for i := 0; i < 1; i++ {

		}
	} else {
		wg := sync.WaitGroup{}
		for i := 0; i < 1; i++ {
			go func() {
				wg.Add(1)
				agent := scenario.Agent{}
				for scen.Tasks.Len() > 0 {
					task := scen.PopTask()
					agent.ExecuteTask(task)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

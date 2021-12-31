package internal

import (
	"elp-go/internal/pathfinding"
	"elp-go/internal/queue"
	"elp-go/internal/world"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

// StartServer Main func when running a server
func StartServer(port int) {
	// listen on the port port
	server, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	for {
		// we accept the incoming connexions on the port port
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		client := NewRemote(conn)
		// Continue to process other client requests
		go func() {
			defer client.Close()
			// A client may make multiple requests with a single connection
			// but not in parallel
			for {
				var scenario Scenario
				err = client.Recv(&scenario)
				if err != nil {
					log.Println(err)
					//log.Fatal(err)
					break
				}
				result := handleRequestSeq(&scenario, pathfinding.NewDijkstra(true, queue.NewLinked))
				log.Printf("sending result : %v", result)
				client.Send(result)
			}
		}()
	}
}

type RequestHandler func(scenario *Scenario, pathfinder pathfinding.Pathfinder) ScenarioResult

func handleRequestSeq(scen *Scenario, pathfinder pathfinding.Pathfinder) ScenarioResult { //test sans goroutine
	result := ScenarioResult{Completed: make([]CompletedTask, len(scen.Tasks))} //make([]CompletedTask, len(scen.Tasks)) : creation of a tab with a length = scen.Tasks

	var idCounter uint32

	for i := 0; i < int(scen.NumAgents); i++ {
		agent := NewAgent(idCounter, world.Pos(0, 0), pathfinder)
		idCounter++
		for j, task := range scen.Tasks {
			task := task.(Task)
			log.Printf("scheduling task %#v on agent %v", task, agent.Id)
			result.Completed[j] = agent.ExecuteTask(scen.World, task)
		}
	}
	return result
}

func handleRequestPar(scen *Scenario, pathfinder pathfinding.Pathfinder) ScenarioResult { //test avec goroutine
	// WaitGroup for all computing goroutines
	agentWg := sync.WaitGroup{}

	// Tasks to be computed, single producer/multi consumers
	tasks := make(chan Task)
	// Task results, multi producers/single consumer
	completed := make(chan CompletedTask, 10)

	var idCounter uint32

	for i := 0; i < int(scen.NumAgents); i++ {
		go func() {
			agentWg.Add(1)
			// Release the lock
			defer agentWg.Done()

			// Initialize a new agent for this coroutine
			agent := NewAgent(atomic.AddUint32(&idCounter, 1), world.Pos(0, 0), pathfinder)
			for t := range tasks {
				log.Printf("scheduling task %#v on agent %v", t, agent.Id)
				completed <- agent.ExecuteTask(scen.World, t)
			}
		}()
	}
	go func() {
		for _, t := range scen.Tasks {
			tasks <- t.(Task)
		}
		close(tasks)

		// Wait for all goroutines to close the output channel
		agentWg.Wait()
		close(completed)
	}()

	compTasks := make([]CompletedTask, len(scen.Tasks))
	i := 0
	for c := range completed {
		compTasks[i] = c
		i++
	}
	return ScenarioResult{Completed: compTasks}
}

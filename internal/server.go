package internal

import (
	"elp-go/internal/pathfinding"
	"elp-go/internal/queue"
	"elp-go/internal/world"
	"log"
	"net"
	"sync"
)

// StartServer Main func when running a server
func StartServer(port int) {
	// listen on the port
	server, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	for {
		// we accept the incoming connexions on the port
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		client := NewRemote(conn)
		log.Printf("New connection from %v\n", client)
		// Continue to process other client requests
		go func() {
			defer client.Close()
			// A client may make multiple requests within a single connection
			// but not in parallel
			for {
				var scenario Scenario
				err = client.Recv(&scenario)
				if err != nil {
					log.Println(err)
					//log.Fatal(err)
					break
				}
				log.Printf("New scenario compute request from %v\n", client)
				result := handleRequestSeq(&scenario, pathfinding.NewAstar(true, pathfinding.EuclideanSq, queue.NewLinked))
				log.Printf("sending result : %v", result)
				client.Send(result)
			}
		}()
	}
}

type RequestHandler func(scenario *Scenario, pathfinder pathfinding.Pathfinder) ScenarioResult

func handleRequestSeq(scen *Scenario, pathfinder pathfinding.Pathfinder) ScenarioResult { //test sans goroutine
	result := ScenarioResult{Completed: make([]CompletedTask, len(scen.Tasks))}

	taskQueue := make(chan Task)
	go func() { // Supply the task queue
		for _, t := range scen.Tasks {
			taskQueue <- t.(Task)
		}
		close(taskQueue)
	}()

	agents := make([]Agent, scen.NumAgents)
	for i := range agents {
		agents[i] = NewAgent(uint(i), world.Pos(0, 0), pathfinder)
	}

	// Dispatch tasks to agent (round-robin)
	index := 0
	for task := range taskQueue {
		//log.Printf("scheduling task %#v on agent %v", task, agents[index%int(scen.NumAgents)].Id)
		result.Completed[index] = agents[index%int(scen.NumAgents)].ExecuteTask(scen.World, task)
		index++
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

	for i := 0; i < scen.NumAgents; i++ {
		i := i // Apparently : Loop variables captured by 'func' literals in 'go' statements might have unexpected values
		go func() {
			agentWg.Add(1)
			// Release the lock
			defer agentWg.Done()

			agent := NewAgent(uint(i), world.Pos(0, 0), pathfinder)
			for t := range tasks {
				//log.Printf("scheduling task %#v on agent %v", t, agent.Id)
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

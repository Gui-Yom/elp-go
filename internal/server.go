package internal

import (
	"elp-go/internal/pathfinding"
	"elp-go/internal/queue"
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
	log.Printf("Server started, listening on %v", server.Addr())
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
					log.Printf("Client %v disconnected\n", client)
					break
				}
				log.Printf("New scenario compute request from %v\n", client)
				result := handleRequestPar(&scenario, pathfinding.NewAstar(true, pathfinding.EuclideanSq, queue.NewLinked))
				//log.Printf("Computed result : %v", result)
				err := client.Send(result)
				if err != nil {
					log.Printf("Client %v disconnected\n", client)
					break
				}
			}
		}()
	}
}

type RequestHandler func(scenario *Scenario, pathfinder pathfinding.Pathfinder) ScenarioResult

// Handle the scenario request sequentially
func handleRequestSeq(scen *Scenario, pathfinder pathfinding.Pathfinder) ScenarioResult { //test sans goroutine
	result := ScenarioResult{Completed: make([]CompletedTask, len(scen.Tasks))}

	taskQueue := make(chan Task)
	go func() { // Supply the task queue
		for _, t := range scen.Tasks {
			taskQueue <- t.(Task)
		}
		close(taskQueue)
	}()

	agents := make([]Agent, len(scen.Agents))
	for i, pos := range scen.Agents {
		agents[i] = NewAgent(uint(i), pos, pathfinder)
	}

	// Dispatch tasks to agent (round-robin)
	index := 0
	for task := range taskQueue {
		//log.Printf("scheduling task %#v on agent %v", task, agents[index%int(scen.Agents)].Id)
		result.Completed[index] = agents[index%len(agents)].ExecuteTask(scen.World, task)
		index++
	}
	return result
}

// Handle scenario request in parallel. Give each agent its goroutine.
func handleRequestPar(scen *Scenario, pathfinder pathfinding.Pathfinder) ScenarioResult { //test avec goroutine
	// WaitGroup for all computing goroutines
	agentWg := sync.WaitGroup{}

	// Tasks to be computed, single producer/multi consumers
	tasks := make(chan Task)
	// Task results, multi producers/single consumer
	completed := make(chan CompletedTask, 10)

	for i := 0; i < len(scen.Agents); i++ {
		i := i // Apparently : Loop variables captured by 'func' literals in 'go' statements might have unexpected values
		go func() {
			agentWg.Add(1)
			// Release the lock
			defer agentWg.Done()

			agent := NewAgent(uint(i), scen.Agents[i], pathfinder)
			for t := range tasks {
				//log.Printf("scheduling task %#v on agent %v", t, agent.Id)
				completed <- agent.ExecuteTask(scen.World, t)
			}
		}()
	}
	go func() {
		// Dispatch all tasks
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
	for c := range completed { // Wait for all results
		compTasks[i] = c
		i++
	}
	return ScenarioResult{Completed: compTasks}
}

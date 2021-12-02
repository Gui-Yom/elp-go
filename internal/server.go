package internal

import (
	"elp-go/internal/pathfinding"
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
				result := handleRequest(&scenario)
				log.Printf("sending result : %v", result)
				client.Send(result)
			}
		}()
	}
}

func handleRequest(scen *Scenario) ScenarioResult {
	log.Printf("recv scenario : %v", scen)

	// Workaround serialization issues
	tasks := make([]Task, len(scen.Tasks))
	for i := range scen.Tasks {
		tasks[i] = scen.Tasks[i].(Task)
	}

	sequential := true
	result := ScenarioResult{Completed: make([]CompletedTask, len(scen.Tasks))}

	if sequential {
		for i := 0; i < int(scen.NumAgents); i++ {
			agent := NewAgent(pathfinding.Pos(0, 0), pathfinding.Dijkstra{Diagonal: scen.DiagonalMovement})
			for _, task := range tasks {
				log.Printf("scheduling task %v on agent %v", task, agent.Id)
				result.Completed = append(result.Completed, agent.ExecuteTask(scen.Carte, task))
			}
		}
	} else {
		agentWg := sync.WaitGroup{}
		pool := NewTaskPool(tasks)
		for i := 0; i < int(scen.NumAgents); i++ {
			go func() {
				agentWg.Add(1)
				// Release the lock
				defer agentWg.Done()

				// Initialize a new agent for this coroutine
				agent := NewAgent(pathfinding.Pos(0, 0), pathfinding.Dijkstra{Diagonal: scen.DiagonalMovement})
				for {
					// Pick a task from the task pool
					task := pool.Get()
					// If no task is available, we just quit
					if task == nil {
						break
					}
					log.Printf("scheduling task %v on agent %v", task, agent.Id)
					result.Completed = append(result.Completed, agent.ExecuteTask(scen.Carte, task))
				}
			}()
		}
		// Wait for all locks
		agentWg.Wait()
	}
	return result
}

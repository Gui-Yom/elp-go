package scenario

import (
	"fmt"
	"log"
)

type Agent struct {
	Id  uint
	Pos Position
}

func (a Agent) ExecuteTask(task Task) {
	switch t := task.(type) {
	case MoveTask:
		log.Printf("t: %v", t)
	default:
		log.Fatalf("Unimplemented task : %v", t)
	}
}

//boucle qui lit la liste de tâches
func findTask(listTasks []Task) {
	if len(listTasks) != 0 {
		//j'essaie de trouver une tâche
		//je la fais
		//je la coche
		//removeElement(listTasks,taskToDo,0)
		//et rebelote
		findTask(listTasks)
	} else {
		fmt.Println("La liste de tâches est vide")
	}
}

func removeElement(s []Task, taskCompleted Task, index int) []Task {
	return append(s[:index], s[index+1:]...)
}

package scenario

import(
	"fmt"
)

type Agent struct {
	//id
	//position
}

//boucle qui lit la liste de tâches
func findTask(listTasks [] Task) {
	if len(listTasks) != 0{
		//j'essaie de trouver une tâche
		//je la fais
		//je la coche
		//et rebelote
		findTask(listTasks)
	}else{
		fmt.Println("La liste de tâches est vide")
	}
}

func removeElement(s []Task, taskCompleted Task, index int) []Task {
	return append(s[:index], s[index+1:]...)
}

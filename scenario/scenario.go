package scenario

type Scenario struct {
	Carte *Carte
	// agents
	Agent *Agent
	// liste de taches
	
	// parametres
}

type Task interface {
	a()
}



type MoveTask struct {
	goal Position
}

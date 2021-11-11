package scenario

type Scenario struct {
	Carte *Carte
	// agents
	Agent *Agent
	// taches
	// parametres
}

type Task interface {
	a()
}

type MoveTask struct {
	goal Position
}

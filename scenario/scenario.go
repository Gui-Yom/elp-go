package scenario

type Scenario struct {
	Carte *Carte
	// agents
	// taches
	// parametres
}

type Task interface {
	a()
}

type MoveTask struct {
	goal Position
}

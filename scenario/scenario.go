package scenario

type Scenario struct {
	Carte *Carte
	// agents
	Agent *Agent
	// liste de taches
	
	// parametres
}

type Task interface {
	howlong() float64
}



type MoveTask struct {
	//goal Position
	name string
}

func (m MoveTask) howlong() float64{
	return 5
}

type fixingBike struct {
	name string
}

func (f fixingBike) howlong() float64{
	return 7
}

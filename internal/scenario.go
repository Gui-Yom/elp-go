package internal

import (
	"elp-go/internal/pathfinding"
	"elp-go/internal/world"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Scenario struct {
	World            *world.World
	DiagonalMovement bool
	Tasks            []interface{}
	Agents           []world.Position
}

func LoadFromFile(filename string) Scenario {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var data map[string]interface{}
	if json.NewDecoder(file).Decode(&data) != nil {
		panic("Can't read scenario from file")
		return Scenario{}
	}

	scen := Scenario{World: mapFromArgs(path.Dir(filename), strings.Split(data["map"].(string), " ")), DiagonalMovement: data["diagonal"].(bool)}

	tasks := data["tasks"]
	switch tasks.(type) {
	case float64:
		scen.Tasks = genTasks(int(tasks.(float64)), scen.World)
	case []interface{}:
		temp := tasks.([]interface{})
		for i, pos := range temp {
			temp2 := pos.(map[string]interface{})
			temp[i] = MoveTask{Goal: world.Pos(int(temp2["x"].(float64)), int(temp2["y"].(float64)))}
		}
		scen.Tasks = temp
	}

	agents := data["agents"]
	switch agents.(type) {
	case float64:
		scen.Agents = make([]world.Position, int(agents.(float64)))
	case []interface{}:
		temp := agents.([]interface{})
		scen.Agents = make([]world.Position, len(temp))
		for i, pos := range temp {
			temp2 := pos.(map[string]interface{})
			scen.Agents[i] = world.Pos(int(temp2["x"].(float64)), int(temp2["y"].(float64)))
		}
	}

	return scen
}

func mapFromArgs(basedir string, args []string) *world.World {
	parseError := func(value, name, type_ string) {
		fmt.Printf("Can't parse %v '%v' as a valid %v", name, value, type_)
		os.Exit(-1)
	}

	var carte *world.World
	argsLen := len(args)
	// Parsing des arguments de création de map
	if argsLen == 0 || args[0] == "rand" {
		var width int = 16
		if argsLen >= 2 {
			if num, err := strconv.ParseInt(args[1], 10, 32); err == nil {
				width = int(num)
			} else {
				parseError(args[1], "width", "int")
			}
		}
		var height int = 16
		if argsLen >= 3 {
			if num, err := strconv.ParseInt(args[2], 10, 32); err == nil {
				height = int(num)
			} else {
				parseError(args[2], "height", "int")
			}
		}
		var fill float32 = 0.2
		if argsLen >= 4 {
			if num, err := strconv.ParseFloat(args[3], 32); err == nil {
				fill = float32(num)
			} else {
				parseError(args[3], "fill", "float")
			}
		}
		var seed int64 = 42
		if argsLen >= 5 {
			if num, err := strconv.ParseInt(args[4], 10, 64); err == nil {
				seed = num
			} else {
				parseError(args[4], "seed", "int")
			}
		}
		carte = world.NewWorldRandom(width, height, fill, seed)
	} else {
		carte = world.NewWorldFromFile(path.Join(basedir, args[0]))
	}
	return carte
}

func genTasks(num int, w *world.World) []interface{} {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tasks := make([]interface{}, num)
	for i := 0; i < len(tasks); i++ {
		pos := world.Pos(rng.Intn(w.Width), rng.Intn(w.Height))
		for ; !w.GetTile(pos).IsTraversable(); pos = world.Pos(rng.Intn(w.Width), rng.Intn(w.Height)) {
		}
		tasks[i] = MoveTask{Goal: pos}
	}
	return tasks
}

type CompletedTask struct {
	AgentId uint
	Path    []world.Position
	Stats   pathfinding.Stats
}

type ScenarioResult struct {
	Completed []CompletedTask
}

func (res ScenarioResult) String() string {
	s := "ScenarioResult\n  Tasks:\n"
	for _, t := range res.Completed {
		s = fmt.Sprintf("%s    %v -> %v (%v)\n", s, t.AgentId, t.Path, t.Stats)
	}
	return s
}

type Task interface {
	Execute(agent *Agent)
}

type MoveTask struct {
	Goal world.Position
}

func (this MoveTask) Execute(agent *Agent) {

}

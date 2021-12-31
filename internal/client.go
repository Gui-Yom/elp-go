package internal

import (
	"bufio"
	"elp-go/internal/world"
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
	"image"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

//fonction qui permet de récupérer un input
func getInput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')

	return strings.TrimSpace(input), err
}

func fillMyList(l []Task) {
	reader := bufio.NewReader(os.Stdin)
	opt, _ := getInput("Choose option (a -add a task, -s save the list): ", reader)

	switch opt {
	case "a":
		fmt.Println("you choose to add a task")
		fillMyList(l)
	case "s":
		fmt.Println("you choose to save the list", l)
	default:
		fmt.Println("that was not a valid option...")
		fillMyList(l)
	}
}

func mapFromArgs(args []string) *world.World {
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
		carte = world.NewWorldFromFile(args[0])
	}
	return carte
}

// StartClient Main func when running a client
func StartClient(addr net.IP, port int, gui bool, connect bool, mapArgs []string) {

	//myListTasks := []internal.Task{}
	//fillMyList(myListTasks)

	carte := mapFromArgs(mapArgs)
	tasks := []interface{}{MoveTask{Goal: world.Pos(16, 16)}}
	scen := Scenario{World: carte, DiagonalMovement: true, NumAgents: 1, Tasks: tasks}

	wait := sync.Mutex{}
	if connect {
		wait.Lock()
		go func() {
			defer wait.Unlock()

			conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: addr, Port: port})
			if err != nil {
				log.Fatal(err)
			}

			client := NewRemote(conn)
			defer client.Close()

			err = client.Send(&scen)
			if err != nil {
				log.Fatal(err)
			}
			var result ScenarioResult
			if err := client.Recv(&result); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("result : %v", result)
		}()
	}

	if gui {
		showScenario(&scen)
	}
	// This effectively waits for the goroutine to end
	wait.Lock()
}

// showScenario Display the given scenario in a window
func showScenario(scen *Scenario) {
	window := app.NewWindow(app.Title("elp-go"), app.Size(unit.Px(720), unit.Px(720)))

	carte := scen.World
	canvas := giocanvas.NewCanvas(720, 720, system.FrameEvent{})

	black := giocanvas.ColorLookup("black")
	white := giocanvas.ColorLookup("white")
	red := giocanvas.ColorLookup("red")

	inputTag := new(bool)

	size := image.Pt(720, 720)
	var mousePos f32.Point

	for e := range window.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			// The window was closed.
			return
		case system.FrameEvent:

			if size != e.Size {
				size = e.Size
				canvas = giocanvas.NewCanvas(float32(size.X), float32(size.Y), e)
			}

			// Process events that arrived between the last frame and this one.
			for _, ev := range e.Queue.Events(inputTag) {
				if x, ok := ev.(pointer.Event); ok {
					switch x.Type {
					case pointer.Move:
						mousePos = x.Position
					}
				}
			}

			// Interested in pointer events
			pointer.InputOp{
				Tag:   inputTag,
				Types: pointer.Move,
			}.Add(canvas.Context.Ops)

			canvas.Background(white)

			tileWidth := float32(size.X) / float32(carte.Width)
			tileHeigth := float32(size.Y) / float32(carte.Height)
			for j := 0; j < carte.Height; j++ {
				for i := 0; i < carte.Width; i++ {
					switch carte.GetTile(world.Pos(i, j)) {
					case world.TILE_EMPTY:
					case world.TILE_WALL:
						canvas.AbsRect(float32(i)*tileWidth, float32(j)*tileHeigth, tileWidth, tileHeigth, black)
					default:
						canvas.AbsRect(float32(i)*tileWidth, float32(j)*tileHeigth, tileWidth, tileHeigth, red)
					}
				}
			}
			canvas.AbsText(1, float32(size.Y-20), 13, fmt.Sprintf("N: %v, Diagonal: %v, mouse %v", scen.NumAgents, scen.DiagonalMovement, mousePos), black)

			// Update the display.
			e.Frame(canvas.Context.Ops)
		}
	}
}

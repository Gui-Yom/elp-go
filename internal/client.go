package internal

import (
	"elp-go/internal/world"
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
	"image"
	"image/color"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

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
func StartClient(addr net.IP, port int, gui bool, connect bool, filename string) {

	scenario := LoadFromFile(filename)

	// Mutex to make the main goroutine wait on the connect
	wait := sync.Mutex{}
	if connect {
		wait.Lock()
		go func() {
			// Unlock when we receive the response
			defer wait.Unlock()

			conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: addr, Port: port})
			if err != nil {
				log.Fatal(err)
			}

			client := NewRemote(conn)
			defer client.Close()

			err = client.Send(&scenario)
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
		// This blocks until we close the window
		showScenario(&scenario)
	}
	// This waits for the connect goroutine to unlock
	wait.Lock()
}

// showScenario Display the given scenario in a gui
func showScenario(scen *Scenario) {
	window := app.NewWindow(app.Title("elp-go"), app.Size(unit.Px(720), unit.Px(720)))

	carte := scen.World
	canvas := giocanvas.NewCanvas(720, 720, system.FrameEvent{})

	black := giocanvas.ColorLookup("black")
	white := giocanvas.ColorLookup("white")
	red := giocanvas.ColorLookup("red")
	sand := color.NRGBA{R: 255, G: 203, B: 107, A: 255}
	belt := color.NRGBA{R: 80, G: 96, B: 88, A: 255}

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
					case world.TILE_SAND:
						canvas.AbsRect(float32(i)*tileWidth, float32(j)*tileHeigth, tileWidth, tileHeigth, sand)
					case world.TILE_CONVEYOR_BELT:
						canvas.AbsRect(float32(i)*tileWidth, float32(j)*tileHeigth, tileWidth, tileHeigth, belt)
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

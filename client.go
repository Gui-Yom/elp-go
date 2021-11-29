package main

import (
	"bufio"
	"elp-go/internal"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/fogleman/gg"
	"image"
	"image/color"
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

func fillMyList(l []internal.Task) {
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

// StartClient Main func when running a client
func StartClient(addr string, port int, gui bool, mapArgs []string) {

	//myListTasks := []internal.Task{}
	//fillMyList(myListTasks)

	parseError := func(value, name, type_ string) {
		fmt.Printf("Can't parse %v '%v' as a valid %v", name, value, type_)
		os.Exit(-1)
	}

	var carte *internal.Carte
	// Parsing des arguments de création de map
	if len(mapArgs) == 0 || mapArgs[0] == "rand" {
		var width int = 16
		if len(mapArgs) >= 2 {
			if num, err := strconv.ParseInt(mapArgs[1], 10, 32); err == nil {
				width = int(num)
			} else {
				parseError(mapArgs[1], "width", "int")
			}
		}
		var height int = 16
		if len(mapArgs) >= 3 {
			if num, err := strconv.ParseInt(mapArgs[2], 10, 32); err == nil {
				height = int(num)
			} else {
				parseError(mapArgs[2], "height", "int")
			}
		}
		var fill float32 = 0.2
		if len(mapArgs) >= 4 {
			if num, err := strconv.ParseFloat(mapArgs[3], 32); err == nil {
				fill = float32(num)
			} else {
				parseError(mapArgs[3], "fill", "float")
			}
		}
		var seed int64 = 42
		if len(mapArgs) >= 5 {
			if num, err := strconv.ParseInt(mapArgs[4], 10, 64); err == nil {
				seed = num
			} else {
				parseError(mapArgs[4], "seed", "int")
			}
		}
		carte = internal.NewMapRandom(width, height, fill, seed)
	} else {
		carte = internal.NewMapFromFile(mapArgs[0])
	}
	tasks := []interface{}{internal.MoveTask{Goal: internal.Pos(16, 16)}}
	scen := internal.Scenario{Carte: carte, DiagonalMovement: true, NumAgents: 1, Tasks: tasks}

	wait := sync.Mutex{}
	wait.Lock()
	go func() {
		defer wait.Unlock()

		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.ParseIP(addr), Port: port})
		if err != nil {
			log.Fatal(err)
		}

		client := NewRemote(conn)
		defer client.Close()

		err = client.Send(&scen)
		if err != nil {
			log.Fatal(err)
		}
		var result internal.ScenarioResult
		if err := client.Recv(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("result : %v", result)
	}()

	if gui {
		showScenario(&scen)
	}
	// This effectively waits for the goroutine to end
	wait.Lock()
}

func showScenario(scen *internal.Scenario) {
	a := app.NewWithID("marais.elp-go")
	//a.SetIcon(theme.SearchIcon())
	w := a.NewWindow("Map display")
	w.SetPadded(false)
	w.SetMaster()

	carte := scen.Carte

	canvas := container.NewMax(canvas.NewRaster(func(w, h int) image.Image {
		ctx := gg.NewContext(w, h)
		tileWidth := float64(w) / float64(carte.Width)
		tileHeigth := float64(h) / float64(carte.Height)
		for j := 0; j < carte.Height; j++ {
			for i := 0; i < carte.Width; i++ {
				ctx.DrawRectangle(float64(i)*tileWidth, float64(j)*tileHeigth, tileWidth, tileHeigth)
				switch carte.GetTile(internal.Pos(i, j)) {
				case internal.TILE_EMPTY:
					ctx.SetColor(color.White)
				case internal.TILE_WALL:
					ctx.SetRGB(0.0, 0.0, 0.0)
				default:
					ctx.SetRGB(1.0, 0.0, 0.0)
				}
				ctx.Fill()
			}
		}
		ctx.SetRGB(0.0, 0.0, 1.0)
		ctx.DrawString(fmt.Sprintf("N: %v, Diagonal: %v", scen.NumAgents, scen.DiagonalMovement), 5, float64(h-10))
		return ctx.Image()
	}))

	w.SetContent(canvas)

	w.Resize(fyne.NewSize(720, 720))
	w.ShowAndRun()
}

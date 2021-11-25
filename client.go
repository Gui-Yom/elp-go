package main

import (
	"bufio"
	"elp-go/scenario"
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
	"strings"
)

//fonction qui permet de récupérer un input
func getInput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')

	return strings.TrimSpace(input), err
}

func fillMyList(l []scenario.Task) {
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
func StartClient(addr string, port int) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.ParseIP(addr), Port: port})
	if err != nil {
		log.Fatal(err)
	}

	client := NewRemote(conn)
	defer client.Close()

	myListTasks := []scenario.Task{}
	fillMyList(myListTasks)
	carte := scenario.NewMapFromFile("map0.map")
	fmt.Printf("%v\n", carte)
	go func() {
		scenario := scenario.Scenario{Carte: carte}

		err = client.Send(&scenario)
		if err != nil {
			log.Fatal(err)
		}
	}()

	showMap(carte)
}

func showMap(carte *scenario.Carte) {
	a := app.NewWithID("marais.elp-go")
	//a.SetIcon(theme.SearchIcon())
	w := a.NewWindow("Map display")
	w.SetPadded(false)
	w.SetMaster()

	canvas := container.NewMax(canvas.NewRaster(func(w, h int) image.Image {
		ctx := gg.NewContext(w, h)
		tileWidth := float64(w) / float64(carte.Width)
		tileHeigth := float64(h) / float64(carte.Height)
		for j := 0; j < carte.Height; j++ {
			for i := 0; i < carte.Width; i++ {
				ctx.DrawRectangle(float64(i)*tileWidth, float64(j)*tileHeigth, tileWidth, tileHeigth)
				switch carte.GetTile(scenario.Pos(i, j)) {
				case scenario.TILE_EMPTY:
					ctx.SetColor(color.White)
				case scenario.TILE_WALL:
					ctx.SetRGB(0.0, 0.0, 0.0)
				default:
					ctx.SetRGB(1.0, 0.0, 0.0)
				}
				ctx.Fill()
			}
		}
		ctx.DrawString("test", 10, 10)
		return ctx.Image()
	}))

	w.SetContent(canvas)

	w.Resize(fyne.NewSize(720, 720))
	w.ShowAndRun()
}

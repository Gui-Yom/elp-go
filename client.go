package main

import (
	"bufio"
	"elp-go/scenario"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

//fonction qui permet de récupérer un input
func getInput(prompt string, r *bufio.Reader) (string,error){
	fmt.Print(prompt)
	input,err := r.ReadString('\n')

	return strings.TrimSpace(input), err
}

func fillMyList(l []scenario.Task){
	reader := bufio.NewReader(os.Stdin)
	opt, _ := getInput("Choose option (a -add a task, -s save the list): ", reader)

	switch opt {
	case "a":
		fmt.Println("you choose to add a task")
		fillMyList(l)
	case "s":
		fmt.Println("you choose to save the list",l)
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
	carte := scenario.ReadMapFromFile("map0.map")
	fmt.Printf("%v\n", carte)
	scenario := scenario.Scenario{Carte: carte}

	err = client.Send(&scenario)
	if err != nil {
		log.Fatal(err)
	}
}

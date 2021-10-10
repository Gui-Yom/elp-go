package main

import (
	"elp-go/client"
	"elp-go/server"
	"flag"
)

func main() {
	// TODO(guillaume) use cobra for arg parsing

	var port int
	flag.IntVar(&port, "p", 32145, "Specify the port to connect or listen to.")
	var startServer bool
	flag.BoolVar(&startServer, "server", false, "Start a startServer.")
	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1", "Host to connect to.")

	flag.Parse()

	if startServer {
		server.Start(port)
	} else {
		client.Start(addr, port)
	}
}

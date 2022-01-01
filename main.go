package main

import (
	"elp-go/internal"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	// CLI Args
	var port int
	flag.IntVar(&port, "p", 32145, "Specify the port to connect or listen to.")
	var startServer bool
	flag.BoolVar(&startServer, "server", false, "Start a server.")
	addr := net.ParseIP("127.0.0.1")
	flag.Func("addr", "Host to connect to.", func(s string) error {
		addr = net.ParseIP(s)
		if addr == nil {
			return errors.New("could not parse addr as an IP")
		}
		return nil
	})
	var nogui bool
	flag.BoolVar(&nogui, "nogui", false, "Disable GUI.")
	var noconnect bool
	flag.BoolVar(&noconnect, "noconnect", false, "Do not connect to remote server.")

	flag.Usage = func() {
		fmt.Println("Usage : elp-go [-args] [scenario file]")
		flag.PrintDefaults()
		fmt.Println("Example usage : \n" +
			"  Start a server :\n" +
			"    $ elp-go -server\n" +
			"  Start a client with a scenario file :\n" +
			"    $ elp-go scenarios/scen0.scen\n" +
			"  Start a client with a random scenario :\n" +
			"    $ elp-go scenarios/rng.scen")
	}

	flag.Parse()

	// Start the server or the client depending on the flags
	if startServer {
		internal.StartServer(port)
	} else {
		if len(flag.Args()) != 1 {
			fmt.Println("Wrong number of arguments")
			os.Exit(-1)
		}
		internal.StartClient(addr, port, !nogui, !noconnect, flag.Args()[0])
	}
}

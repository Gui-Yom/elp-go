package main

import (
	"elp-go/internal"
	"errors"
	"flag"
	"fmt"
	"net"
)

func main() {
	// CLI Args
	var port int
	flag.IntVar(&port, "p", 32145, "Specify the port to connect or listen to.")
	var startServer bool
	flag.BoolVar(&startServer, "server", false, "Start a server.")
	var addr net.IP
	addr = net.ParseIP("127.0.0.1")
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
		fmt.Println("Usage : elp-go [-args] [map file | rand] [width] [height] [fill] [seed]")
		flag.PrintDefaults()
		fmt.Println("Example usage : \n" +
			"  Start a server :\n" +
			"    $ elp-go -server\n" +
			"  Start a client with a map file :\n" +
			"    $ elp-go map.map\n" +
			"  Start a client with a randomly generated map :\n" +
			"    $ elp-go rand 100 100 0.1 42")
	}

	flag.Parse()

	if startServer {
		internal.StartServer(port)
	} else {
		internal.StartClient(addr, port, !nogui, !noconnect, flag.Args())
	}
}

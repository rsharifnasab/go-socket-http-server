package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

const DEFAULT_PORT int = 80

const FLAG_PORT_HELP string = "the port to listen for requests"

const CANNOT_OPEN_PORT_MSG string = `cannot connect to specified port
 consider running executable with super user permission
 hint: opening ports <= 1024 typically need more permission`

var flagPort *int = flag.Int("port", DEFAULT_PORT, FLAG_PORT_HELP)

func main() {

	flag.Parse()

	address := fmt.Sprintf(":%d", *flagPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(CANNOT_OPEN_PORT_MSG)
	}
	log.Printf("server is running on " + address)
	_ = listener
}

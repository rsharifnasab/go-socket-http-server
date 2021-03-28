package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

const DEFAULT_PORT int = 80

const FLAG_PORT_HELP string = "the port to listen for requests"

const CANNOT_OPEN_PORT_MSG string = `cannot connect to specified port
 hint: consider running executable with superuser permission
 because opening ports <= 1024 typically need more permission`

var flagPort *int = flag.Int("port", DEFAULT_PORT, FLAG_PORT_HELP)

func main() {

	flag.Parse()

	address := fmt.Sprintf(":%d", *flagPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(CANNOT_OPEN_PORT_MSG)
	}
	defer listener.Close()
	log.Printf("server is running on " + address)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("error in connecting [%s]", err)
		}

		go connectionHandler(connection)
		// go routine: handle it concurrently
	}
}

func connectionHandler(conn net.Conn) {
	defer conn.Close()
	log.Printf("new connection from [%s]", conn.LocalAddr())
	writer := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(conn)

	writer.WriteString("salam chetori?\n")
	writer.Flush()

	for scanner.Scan() {
		scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		// error: TODO
		writer.WriteString("nashod")
	}
	log.Printf("connection closed")
}

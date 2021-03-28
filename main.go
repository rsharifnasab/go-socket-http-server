package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

const DEFAULT_PORT = 80
const ALTERNATE_PORT = 8080

func checkAvailablePort(port int) error {
	portAddStr := fmt.Sprintf(":%d", port)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", portAddStr)
	if err != nil {
		log.Fatal("there was a problem in using 127.0.0.1 or port number")
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		// cannot open port
		return err
	} else {
		listener.Close()
		return nil
	}
}

func main() {
	var port int = DEFAULT_PORT
	if checkAvailablePort(DEFAULT_PORT) != nil {
		log.Printf("cannot open %d port :(\ntry running with sudo access\n", DEFAULT_PORT)
		if checkAvailablePort(ALTERNATE_PORT) != nil {
			log.Fatalf("cannot open %d port either, exiting\n", ALTERNATE_PORT)
		} else {
			port = ALTERNATE_PORT
			log.Printf("%d port was ok, using that\n", ALTERNATE_PORT)
		}
	} else {
		log.Printf("using port %d", port)
	}
	portStr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("/ called")
		fmt.Println(req.URL)
	})
	http.HandleFunc("/author", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "author: roozbeh sharifnasab\n")
		log.Println("author called")
	})
	log.Fatal(http.ListenAndServe(portStr, nil))
}

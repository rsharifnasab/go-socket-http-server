package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

const DEFAULT_PORT int = 80

const FLAG_PORT_HELP string = "the port to listen for requests"

const CANNOT_OPEN_PORT_MSG string = `cannot connect to specified port
 hint: consider running executable with superuser permission
 because opening ports <= 1024 typically need more permission`

var HEADER_REGEX_1 = regexp.MustCompile(
	`^(GET|POST) (/[\w\.]*) HTTP/\d\.\d$`)

var flagPort *int = flag.Int("port", DEFAULT_PORT, FLAG_PORT_HELP)

type (
	HttpRequest struct {
		Method string
		Path   string
	}

	HttpResponse struct {
		Status        int
		Date          string
		ContentLength int
		Data          []byte
	}
)

func main() {

	flag.Parse()

	address := fmt.Sprintf(":%d", *flagPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(CANNOT_OPEN_PORT_MSG)
	}
	defer listener.Close()

	log.Printf("Server is running on " + address)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("error in connecting [%s]", err)
		} else {
			go connectionHandler(connection)
			// go routine: handle it concurrently
		}
	}
}

func CreateRequest(scanner *bufio.Scanner) (*HttpRequest, error) {
	req := &HttpRequest{}

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	firstLine := scanner.Text()
	res := HEADER_REGEX_1.FindStringSubmatch(firstLine)
	if len(res) < 1 {
		return nil, fmt.Errorf(
			"[%s] doesn't belong to a valid http request", firstLine)
	}
	req.Method = res[1]
	req.Path = res[2]
	fmt.Println(req)

	for scanner.Scan() {

		_ = strings.Split(strings.TrimSpace(scanner.Text()), " ")
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return req, nil
}

func writeResponse(response HttpResponse, writer *bufio.Writer) error {
	writer.WriteString("salam chetori?\n")
	writer.Flush()

	return nil
}

func serverLogic(req *HttpRequest) HttpResponse {
	response := HttpResponse{Status: 200}

	return response
}

func connectionHandler(conn net.Conn) {
	defer conn.Close()
	log.Printf("new connection from [%s]", conn.LocalAddr())

	request, err := CreateRequest(bufio.NewScanner(conn))
	if err != nil {
		return // 500
	}

	response := serverLogic(request)

	writeResponse(response, bufio.NewWriter(conn))

	log.Printf("connection closed")
}

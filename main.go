package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

const DEFAULT_PORT int = 80

const FLAG_PORT_HELP string = "the port to listen for requests"

const CANNOT_OPEN_PORT_MSG string = `cannot connect to specified port
 hint: consider running executable with superuser permission
 because opening ports <= 1024 typically need more permission`

var HEADER_REGEX_1 = regexp.MustCompile(
	`^(GET|POST) (/[\w\.]*) HTTP/\d\.\d$`)

var flagPort *int = flag.Int("port", DEFAULT_PORT, FLAG_PORT_HELP)

const EXAMPLE_TIME_FORMAT = `Mon, 29 Mar 2021 13:21:49 GMT`

const NOT_FOUND_MSG = `404 Not Found`

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

	for scanner.Scan() {
		_ = strings.TrimSpace(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return req, nil
}

func writeResponse(response *HttpResponse, writer *bufio.Writer) error {
	writer.WriteString("salam chetori?\n")
	writer.Flush()

	return nil
}

func serverLogic(req *HttpRequest) *HttpResponse {
	response := &HttpResponse{Status: 200}

	return response
}

func emptyResponse() *HttpResponse {
	return &HttpResponse{
		ContentLength: 0,
		Status:        200,
		Date:          time.Now().Format(EXAMPLE_TIME_FORMAT),
	}
}

func createError400() *HttpResponse {
	resp := emptyResponse()
	resp.Status = 400
	return resp
}

func createError404() *HttpResponse {
	resp := emptyResponse()
	resp.Status = 400
	resp.Data = []byte(NOT_FOUND_MSG)
	return resp
}
func connectionHandler(conn net.Conn) {
	defer conn.Close()
	log.Printf("new connection from [%s]", conn.LocalAddr())

	writer := bufio.NewWriter(conn)

	request, err := CreateRequest(bufio.NewScanner(conn))
	if err != nil {
		_ = writeResponse(createError400(), writer)
		// we can't handle that error :(
		return
	}

	response := serverLogic(request)

	writeResponse(response, writer)

	log.Printf("connection from [%s] closed", conn.LocalAddr())
}

func extensionToMime(fileExt string) (string, error) {
	switch fileExt {
	case ".html":
		return "text/html", nil
	case ".htm":
		return "text/html", nil
	case ".js":
		return "application/javascript", nil
	case ".json":
		return "application/json", nil
	case ".xml":
		return "application/xml", nil
	case ".zip":
		return "application/zip", nil
	case ".wma":
		return "audio/x-ms-wma", nil
	case ".txt":
		return "text/plain", nil
	case ".ttf":
		return "applicatcase ion/x-font-ttf", nil
	case ".tex":
		return "application/x-tex", nil
	case ".sh":
		return "application/x-sh", nil
	case ".py":
		return "text/x-python", nil
	case ".png":
		return "image/png", nil
	case ".pdf":
		return "application/pdf", nil
	case ".mpeg":
		return "video/mpeg", nil
	case ".mpa":
		return "video/mpeg", nil
	case ".mp4":
		return "video/mp4", nil
	case ".mp3":
		return "audio/mpeg", nil
	case ".log":
		return "text/plain", nil
	case ".jpg":
		return "image/jpeg", nil
	case ".jpeg":
		return "image/jpeg", nil
	case ".java":
		return "text/x-java-source", nil
	case ".jar":
		return "application/java-archive", nil
	case ".gif":
		return "image/gif", nil
	case ".cpp":
		return "text/x-c", nil
	case ".bmp":
		return "image/bmp", nil
	case ".avi":
		return "video/x-msvideo", nil
	case ".mkv":
		return "video/x-matroska", nil

	default:
		return "application/octet-stream",
			fmt.Errorf("[%s] is not a valid extension", fileExt)
	}

}

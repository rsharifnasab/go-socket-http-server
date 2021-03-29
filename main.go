package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

const DEFAULT_PORT int = 80

const ROOT = "./static/"

const FLAG_PORT_HELP string = "the port to listen for requests"

const CANNOT_OPEN_PORT_MSG string = `cannot connect to specified port
 hint: consider running executable with superuser permission
 because opening ports <= 1024 typically need more permission`

var HEADER_REGEX_1 = regexp.MustCompile(
	`^(GET|POST) (/[\w\./]*) HTTP/\d\.\d$`)

var flagPort *int = flag.Int("port", DEFAULT_PORT, FLAG_PORT_HELP)

const NOT_FOUND_MSG = "404 Not Found\r\n"

type (
	HttpRequest struct {
		Method string
		Path   string
	}

	HttpResponse struct {
		Status        int
		Date          time.Time
		ContentLength int64
		Data          []byte
		Reader        io.Reader
		LastModified  time.Time
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

	// read headers
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if len(text) == 0 {
			// reached empty line after headers
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return req, nil
}

func httpCodeToStatus(code int) string {
	switch code {
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"

	default:
		log.Fatalf("unknown status code [%d]", code)
		return "501 Not Implemented"
	}
}

func writeResponse(res *HttpResponse, writer *bufio.Writer) error {
	var format = fmt.Sprintf
	firstLine := format("HTTP/1.1 %d %s\r\n",
		res.Status, httpCodeToStatus(res.Status))

	headers := format("Data: %s\r\n", res.Date.Format(time.RFC1123Z))
	headers += "Server: RoozbehsServer/1.0\r\n"
	headers += format("Last-Modified: %s\r\n",
		res.LastModified.Format(time.RFC1123Z))
	headers += "Connection: close\r\n"
	headers += format("Content-Length: %d\r\n", res.ContentLength)

	_, err := writer.WriteString(firstLine + headers + "\r\n")
	if err != nil {
		return err
	}
	if res.Reader != nil {
		lenn, err := writer.ReadFrom(res.Reader)
		if err != nil {
			return err
		}
		if lenn != res.ContentLength {
			return errors.New("sent file length wasn't same as file length")
		}
	} else {
		writer.Write(res.Data)
	}

	writer.Flush()
	return nil
}

func createResponse(req *HttpRequest) (*HttpResponse, *os.File) {
	response := emptyResponse()
	response.Status = 200

	file, err := os.Open(ROOT + req.Path)
	if err != nil {
		log.Printf("404 [%s]", err)
		return createError404(), nil
	}
	response.Reader = bufio.NewReader(file)
	fileStat, err := file.Stat()
	if err != nil {
		log.Printf("404 [%s] (err in reading file stat)", err)
		return createError404(), nil
	}
	if fileStat.IsDir() {

	}
	response.ContentLength = fileStat.Size()
	response.LastModified = fileStat.ModTime()

	return response, file
}

func emptyResponse() *HttpResponse {
	return &HttpResponse{
		ContentLength: 0,
		Status:        200,
		Date:          time.Now(),
	}
}

func createError400() *HttpResponse {
	resp := emptyResponse()
	resp.Status = 400
	resp.Data = []byte("error 400, bad request\r\n")
	resp.ContentLength = int64(len(resp.Data))
	return resp
}

func createError404() *HttpResponse {
	resp := emptyResponse()
	resp.Status = 400
	resp.Data = []byte(NOT_FOUND_MSG)
	resp.ContentLength = int64(len(resp.Data))
	return resp
}
func connectionHandler(conn net.Conn) {
	log.Printf("new connection from [%s]", conn.LocalAddr())
	defer conn.Close()

	writer := bufio.NewWriter(conn)

	request, err := CreateRequest(bufio.NewScanner(conn))
	if err != nil {
		log.Printf("400 : %s", err)
		err = writeResponse(createError400(), writer)
		if err != nil {
			log.Printf("cannot send 400 because [%s]", err)
		}
		return
	}

	response, file := createResponse(request)
	defer file.Close()

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

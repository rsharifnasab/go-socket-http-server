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
	"path/filepath"
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
	`^(GET) (/[\w\./]*) HTTP/\d\.\d$`)

var flagPort *int = flag.Int("port", DEFAULT_PORT, FLAG_PORT_HELP)

const NOT_FOUND_MSG = "404 Not Found\r\n"
const BAD_REQ_MSG = "400 Bad Request\r\n"

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
		Mime          string
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
	if firstLine == "" {
		return nil, errors.New("first line was empty")
		//return CreateRequest(scanner)
	}
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

	headers := format("Date: %s\r\n", res.Date.Format(time.RFC1123Z))
	headers += "Server: RoozbehAwesomeServer/1.0\r\n"

	if res.Reader != nil {
		headers += format("Last-Modified: %s\r\n",
			res.LastModified.Format(time.RFC1123Z))
	}
	headers += format("Content-Type: %s;\r\n", res.Mime)
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
		req.Path += "/index.html"
		return createResponse(req)
	}
	response.ContentLength = fileStat.Size()
	response.LastModified = fileStat.ModTime()
	response.Mime = extensionToMime(filepath.Ext(req.Path))

	return response, file
}

func emptyResponse() *HttpResponse {
	return &HttpResponse{
		ContentLength: 0,
		Status:        200,
		Date:          time.Now(),
		Mime:          extensionToMime(".txt"),
	}
}

func createError400() *HttpResponse {
	resp := emptyResponse()
	resp.Status = 400
	resp.Data = []byte(BAD_REQ_MSG)
	resp.ContentLength = int64(len(resp.Data))
	return resp
}

func createError404() *HttpResponse {
	resp := emptyResponse()
	resp.Status = 404
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

func extensionToMime(fileExt string) string {
	switch fileExt {
	case ".html":
		return "text/html"
	case ".htm":
		return "text/html"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".zip":
		return "application/zip"
	case ".wma":
		return "audio/x-ms-wma"
	case ".txt":
		return "text/plain"
	case ".ttf":
		return "applicatcase ion/x-font-ttf"
	case ".tex":
		return "application/x-tex"
	case ".sh":
		return "application/x-sh"
	case ".py":
		return "text/x-python"
	case ".png":
		return "image/png"
	case ".pdf":
		return "application/pdf"
	case ".mpeg":
		return "video/mpeg"
	case ".mpa":
		return "video/mpeg"
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	case ".log":
		return "text/plain"
	case ".jpg":
		return "image/jpeg"
	case ".jpeg":
		return "image/jpeg"
	case ".java":
		return "text/x-java-source"
	case ".jar":
		return "application/java-archive"
	case ".gif":
		return "image/gif"
	case ".cpp":
		return "text/x-c"
	case ".bmp":
		return "image/bmp"
	case ".avi":
		return "video/x-msvideo"
	case ".mkv":
		return "video/x-matroska"
	case ".ico":
		return "image/x-icon"

	default:
		return "application/octet-stream"
	}

}

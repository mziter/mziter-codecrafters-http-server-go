package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}
		go handleConn(c)
	}

}

func handleConn(c net.Conn) {
	defer c.Close()

	input := make([]byte, 1024)
	_, err := c.Read(input)
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		os.Exit(1)
	}

	req, err := http.ParseRequest(input)
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf("REQ: %v\n", req)
	if req.Path == "/" {
		http.WriteResponse(c, http.StatusCodeOK, http.StatusDescriptionOK, nil, nil)
	} else if strings.HasPrefix(req.Path, "/echo/") {
		echoString := strings.TrimPrefix(req.Path, "/echo/")
		echo(c, echoString)
	} else if strings.HasPrefix(req.Path, "/user-agent") {
		userAgent(c, req.Headers)
	} else if strings.HasPrefix(req.Path, "/files/") {
		filename := strings.TrimPrefix(req.Path, "/files/")
		sendFile(c, filename)
	} else {
		http.WriteResponse(c, http.StatusCodeNotFound, http.StatusDescriptionNotFound, nil, nil)
	}
}

func echo(c net.Conn, echoString string) {
	headers := map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprint(len(echoString)),
	}
	http.WriteResponse(c, http.StatusCodeOK, http.StatusDescriptionOK, headers, []byte(echoString))
}

func userAgent(c net.Conn, headers map[string]string) {
	userAgent, found := headers["User-Agent"]
	if !found {
		http.WriteResponse(c, http.StatusCodeBadRequest, http.StatusDescriptionBadRequest, nil, nil)
	}
	respHeaders := map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprint(len(userAgent)),
	}
	http.WriteResponse(c, http.StatusCodeOK, http.StatusDescriptionOK, respHeaders, []byte(userAgent))
}

func sendFile(c net.Conn, fileName string) {
	dir := os.Args[2]
	fmt.Printf("received args: %v", dir)
	contents, err := os.ReadFile(path.Join(dir, fileName))
	if err != nil {
		fmt.Printf("file read error: %s", err.Error())
		http.WriteResponse(c, http.StatusCodeInternalServiceError, http.StatusDescriptionInternalServiceError, nil, nil)
	}
	respHeaders := map[string]string{
		"Content-Type":   "application/octet-stream",
		"Content-Length": fmt.Sprint(len(contents)),
	}
	http.WriteResponse(c, http.StatusCodeOK, http.StatusDescriptionOK, respHeaders, contents)
}

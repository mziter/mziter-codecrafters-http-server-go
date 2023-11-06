package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	c, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	defer c.Close()

	input := make([]byte, 1024)
	_, err = c.Read(input)
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
	}
	c.Close()
}

func echo(c net.Conn, echoString string) {
	headers := map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprint(len(echoString)),
	}
	http.WriteResponse(c, http.StatusCodeOK, http.StatusDescriptionOK, headers, []byte(echoString))
}

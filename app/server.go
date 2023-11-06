package main

import (
	"fmt"
	"net"
	"os"

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
		http.WriteResponse(c, http.StatusCodeOK, http.StatusDescriptionOK)
	} else {
		http.WriteResponse(c, http.StatusCodeNotFound, http.StatusDescriptionNotFound)
	}
	c.Close()
}

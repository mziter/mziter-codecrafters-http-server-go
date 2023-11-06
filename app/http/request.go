package http

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
}

func ParseRequest(request []byte) (Request, error) {
	r := bytes.NewReader(request)
	s := bufio.NewScanner(r)

	var req Request

	// parse http start line
	firstLineExists := s.Scan()
	if !firstLineExists {
		return req, errors.New("couldn't detect start line of http request")
	}

	startLine := s.Bytes()
	words := bytes.Split(startLine, []byte(" "))
	if len(words) != 3 {
		return req, fmt.Errorf("expected start line to be 3 tokens, but was %d", len(words))
	}
	req.Method = string(words[0])
	req.Path = string(words[1])
	req.Version = string(words[2])

	// parse http headers
	req.Headers = make(map[string]string)
	for s.Scan() {
		line := s.Bytes()
		if bytes.Equal(line, []byte("")) {
			break
		}
		tokens := bytes.Split(line, []byte(": "))
		if len(tokens) != 2 {
			return req, fmt.Errorf("expected request header line to be in k: v format: %s", line)
		}
		req.Headers[string(tokens[0])] = string(tokens[1])
	}

	return req, nil
}

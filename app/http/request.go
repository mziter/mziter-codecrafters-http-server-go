package http

import (
	"bufio"
	"bytes"
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
	firstLine := true
	for s.Scan() {
		line := s.Bytes()
		if firstLine {
			words := bytes.Split(line, []byte(" "))
			if len(words) != 3 {
				return req, fmt.Errorf("expected start line to be 3 tokens, but was %d", len(words))
			}
			req.Method = string(words[0])
			req.Path = string(words[1])
			req.Version = string(words[2])
			firstLine = false
			continue
		}
	}

	return req, nil
}

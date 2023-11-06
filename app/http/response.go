package http

import (
	"fmt"
	"io"
	"strings"
)

const (
	StatusCodeOK                          = 200
	StatusCodeBadRequest                  = 400
	StatusCodeNotFound                    = 404
	StatusCodeInternalServiceError        = 500
	StatusDescriptionOK                   = "OK"
	StatusDescriptionBadRequest           = "Bad Request"
	StatusDescriptionNotFound             = "Not Found"
	StatusDescriptionInternalServiceError = "Internal Service Error"
)

func WriteResponse(w io.Writer, statusCode int, statusDescription string, headers map[string]string, body []byte) {
	var out strings.Builder
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusDescription)
	out.WriteString(statusLine)
	for k, v := range headers {
		out.WriteString(k)
		out.WriteString(": ")
		out.WriteString(v)
		out.WriteString("\r\n")
	}
	out.WriteString("\r\n")
	out.Write(body)
	out.WriteString("\r\n\r\n")
	w.Write([]byte(out.String()))
}

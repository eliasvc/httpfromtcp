package request

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
)

var HTTPMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

// HTTP request, composed by a:
// request-line: this is a start-line but specificially for client requests
// fields/headers
// body
type Request struct {
	RequestLine RequestLine
}

// HTTP start-line, which is called a request-line for client requests,
// is composed of an HTTP method (i.e "GET"), a target (i.e "/"), and
// an HTTP version (i.e "HTTP/1.1")
type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	in, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	request_string := string(in)
	sections := strings.Split(request_string, "\r\n")
	sections_len := len(sections)

	request := Request{}
	log.Println(request_string)
	log.Println(sections_len)
	if sections_len >= 1 {
		requestLine := sections[0]
		parts := strings.Split(requestLine, " ")
		// request-line is composed of Method, Request Target and HTTP Version
		if len(parts) != 3 {
			return nil, fmt.Errorf("Error: Malformed request-line: '%s'", requestLine)
		}

		method := parts[0]
		if !slices.Contains(HTTPMethods, method) {
			return nil, fmt.Errorf("Error: Incorrect HTTP Method: %s", method)
		}

		path := parts[1]
		if !strings.HasPrefix(path, "/") {
			return nil, fmt.Errorf("Error: Malformed Request Target: %s", path)
		}

		httpVersion := parts[2]
		version_parts := strings.Split(parts[2], "/")
		if version_parts[0] != "HTTP" || version_parts[1] != "1.1" {
			return nil, fmt.Errorf("Error: Malformed HTTP version: %s", httpVersion)
		}

		request.RequestLine.Method = method
		request.RequestLine.RequestTarget = path
		request.RequestLine.HttpVersion = version_parts[1]

		return &request, nil
	}

	return nil, fmt.Errorf("Error: Malformed Request: %s", request_string)

}

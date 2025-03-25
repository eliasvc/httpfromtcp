package request

import (
	"io"
)

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
	return nil, nil
}

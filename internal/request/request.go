package request

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"
)

const crlf = "\r\n"

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

	requestLine, err := parseRequestLine(in)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	// The first line that ends in CRLF is your request-line
	crlfIndex := bytes.Index(data, []byte(crlf))
	if crlfIndex == -1 {
		return nil, fmt.Errorf("could not find CRLF in request-line")
	}
	requestLineText := string(data[:crlfIndex])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err
	}
	return requestLine, nil
}

func requestLineFromString(requestLineText string) (*RequestLine, error) {
	parts := strings.Fields(requestLineText)
	// request-line is composed of Method, Request Target and HTTP Version. Anything else
	// is a mistake
	if len(parts) != 3 {
		return nil, fmt.Errorf("Error: Malformed request-line: '%s'", requestLineText)
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

	return &RequestLine{
		Method:        method,
		RequestTarget: path,
		HttpVersion:   version_parts[1],
	}, nil

}

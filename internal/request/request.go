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

// Meant for tracking that state of a Request
// 1: initialized
// 0: done
type State int

// HTTP request, composed by a:
// request-line: this is a start-line but specificially for client requests
// fields/headers
// body
type Request struct {
	RequestLine RequestLine
	State       State
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
	buf := make([]byte, 8)
	// requestLineBuf is the container of all data until crlf is found, which marks
	// the request-line
	var requestLineBuf []byte
	var request Request
	request.State = 1
	for {
		// Try to readh len(buf) data into buf:= make([]byte, 8, 8)
		bytes_read, reader_err := io.ReadFull(reader, buf)
		fmt.Printf("DEBUG: buf: '%v'\n", buf)
		// Concatenate previous request-line contents with buf
		requestLineBuf = slices.Concat(requestLineBuf, buf[:bytes_read])
		fmt.Printf("DEBUG: requestLineBuf: '%v'\n", []byte(requestLineBuf))
		fmt.Printf("DEBUG: %s\n", requestLineBuf)
		// Try to parse what we have to see if we finally have a complete request-line
		nParsed, parser_err := request.parse(requestLineBuf)
		// No crlf found yet, keep reading
		if nParsed == 0 && parser_err == nil && reader_err == nil {
			continue
		}
		// request-line found
		if nParsed > 0 && parser_err == nil {
			break
		}
		// Some error found while parsing
		if nParsed > 0 && parser_err != nil {
			return nil, parser_err
		}
		// Reached a reader error without a succesful request-line parsing
		if reader_err == io.EOF {
			break
		}
	}
	fmt.Println(request)
	return &request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State < 0 || r.State > 1 {
		return 0, fmt.Errorf("Error: unknown state: %d", r.State)
	}

	if r.State == 0 {
		return 0, fmt.Errorf("Error: uninitialized parser")
	}

	requestLine, n, err := parseRequestLine(data)
	if err == nil && n == 0 {
		// If no bytes were processed and there was no error, then we just need more data
		// Signal this by passing a nil error
		return n, nil
	}

	if err != nil {
		return n, err
	}

	r.RequestLine = *requestLine
	// We're done, so the parser can be updated appropriately
	r.State = 0
	return n, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	// The first line that ends in CRLF is your request-line
	crlfIndex := bytes.Index(data, []byte(crlf))
	if crlfIndex == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:crlfIndex])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, crlfIndex, err
	}
	return requestLine, crlfIndex, nil
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
		fmt.Printf("RequestLineText in requestLineFromString: '%v'\n", requestLineText)
		fmt.Printf("'%v'\n", []byte(method))
		fmt.Println("HTTPMethods: ", HTTPMethods)
		fmt.Println(slices.Contains(HTTPMethods, "GET"))
		return nil, fmt.Errorf("Error: Incorrect HTTP Method: '%v'", method)
	}

	path := parts[1]
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("Error: Malformed Request Target: '%s'", path)
	}

	httpVersion := parts[2]
	version_parts := strings.Split(parts[2], "/")
	if version_parts[0] != "HTTP" || version_parts[1] != "1.1" {
		return nil, fmt.Errorf("Error: Malformed HTTP version: '%s'", httpVersion)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: path,
		HttpVersion:   version_parts[1],
	}, nil

}

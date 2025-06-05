package request

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"httpfromtcp/internal/headers"
)

type parserState int

const (
	Done parserState = iota
	Initialized
)

type requestHeaderParserState int

const (
	HeaderDone requestHeaderParserState = iota
	HeaderInitialized
)

// Enum tracks the state of the parser:
// - 0: done
// - 1: initialized
type Request struct {
	RequestLine       RequestLine
	Headers           headers.Headers
	pState            parserState
	headerParserState requestHeaderParserState
}

func (r *Request) parse(data []byte) (int, error) {
	if 
	n, requestLine, err := parseRequestLine(data)
	// parser didn't find CRLF
	if err == nil && n == 0 {
		r.pState = Initialized
		return n, nil
	}
	if err != nil {
		return 0, err
	}

	r.RequestLine = requestLine
	r.pState = Done

	n, h, err := parseHeaders(data[n:])
	return n, nil
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{}
	buf := make([]byte, 0, 512)
	for {
		n, err := reader.Read(buf[len(buf):cap(buf)])
		buf = buf[:len(buf)+n]
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("Error reading request: %q\n", err.Error())
		}

		_, err = request.parse(buf)
		if err != nil {
			return nil, fmt.Errorf("Error parsing request-line: %q\n", err)
		}
		if request.pState == Done && request.headerParserState == HeaderDone {
			return &request, nil
		}

		if len(buf) == cap(buf) {
			// Use the io.ReadAll capacity extension trick
			buf = append(buf, 0)[:len(buf)]
		}
	}
}

// parseRequestLine parses the request line from the HTTP request contained in inputRequest.
// If CRLF is not found, it is assummed that we need more data and so it returns err == nil.
// A successful call returns the number of bytes processed, a RequestLine and err == nil.
func parseRequestLine(inputRequest []byte) (int, RequestLine, error) {
	var requestLine RequestLine
	validMethod := regexp.MustCompile(`[A-Z]+`)
	validHTTPVersion := "HTTP/1.1"

	inputRequestText := string(inputRequest)
	crlfIndex := strings.Index(inputRequestText, "\r\n")
	if crlfIndex == -1 {
		return 0, requestLine, nil
	}
	requestLineText := inputRequestText[:crlfIndex]

	parts := strings.Split(requestLineText, " ")
	// request-line is composed of "method SP request-target SP HTTP-version"
	// where SP is single space. So it has to have only 3 parts if devided by
	// space
	if len(parts) != 3 {
		return crlfIndex, requestLine, fmt.Errorf("Malformed request-line: %q\n", requestLineText)
	}

	method, target, version := parts[0], parts[1], parts[2]

	if !validMethod.MatchString(method) {
		return crlfIndex, requestLine, fmt.Errorf("Malformed method in request-line: %q\n", requestLineText)
	}

	if validHTTPVersion != version {
		return crlfIndex, requestLine, fmt.Errorf("Malformed HTTP version in request-line: %q\n", requestLineText)
	}
	versionParts := strings.Split(version, "/")
	httpVersion := versionParts[1]

	requestLine.Method = method
	requestLine.RequestTarget = target
	requestLine.HttpVersion = httpVersion

	return crlfIndex, requestLine, nil
}

// parseHeaders parses all headers in the HTTP request contained in data
// A successful run will return n bytes consumed, a Headers with the parsed headers,
// and nil for error
func parseHeaders(data []byte) (int, headers.Headers, error) {
	h := headers.NewHeaders()
	done := false
	totalParsed := 0

	for !done {
		n, done, err := h.parse(data)
			
	}
}

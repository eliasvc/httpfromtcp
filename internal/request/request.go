package request

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"httpfromtcp/internal/headers"
)

type parserState int

const (
	Initialized parserState = iota
	Done
)

type requestHeaderParserState int

const (
	HeaderInitialized requestHeaderParserState = iota
	HeaderDone
)

// Enum tracks the state of the parser:
// - 0: done
// - 1: initialized
type Request struct {
	RequestLine       RequestLine
	Headers           headers.Headers
	pState            parserState
	headerParserState requestHeaderParserState
	totalBytesParsed  int
}

const crlf = "\r\n"

func (r *Request) parse(data []byte) (int, error) {
	if r.pState == Initialized {
		n, requestLine, err := parseRequestLine(data)
		// parser didn't find CRLF
		if err == nil && n == 0 {
			return 0, nil
		}
		if err != nil {
			return 0, err
		}
		r.RequestLine = requestLine
		r.pState = Done
		r.totalBytesParsed += n
	}

	if r.pState == Done && r.headerParserState == HeaderInitialized {
		for i := 0; r.headerParserState != HeaderDone; i++ {
			nParsed, done, err := r.Headers.Parse(data[r.totalBytesParsed:])
			if err != nil {
				return 0, err
			}
			// Need to read more data
			if nParsed == 0 && !done {
				return 0, nil
			}

			if done {
				r.headerParserState = HeaderDone
				return r.totalBytesParsed, nil
			}

			r.totalBytesParsed += nParsed
		}

	}

	return r.totalBytesParsed, nil
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{
		Headers: headers.NewHeaders(),
	}
	buf := make([]byte, 0, 512)
	fmt.Println(">> Reading from stream")
	for i := 0; ; i++ {
		//if i == 350 {
		//	fmt.Println(">> FORCED BREAK ON ReadFromReader")
		//	fmt.Printf("buf: %q", buf)
		//	fmt.Printf("totalBytesParsed: %q", request.totalBytesParsed)
		//	return nil, fmt.Errorf("Parsign logic Error\n")
		//}
		n, readErr := reader.Read(buf[len(buf):cap(buf)])
		buf = buf[:len(buf)+n]
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return nil, fmt.Errorf("Error reading request: %q\n", readErr.Error())
		}

		n, parseErr := request.parse(buf)
		if parseErr != nil {
			return nil, fmt.Errorf("Error parsing request-line: %q\n", parseErr)
		}
		if request.pState == Done && request.headerParserState == HeaderDone {
			return &request, nil
		}
		// If EOF was reached and the parsers are not done, something is wrong with the request
		if errors.Is(readErr, io.EOF) {
			return nil, fmt.Errorf("Malformed Request: %q\n", buf)
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
		return crlfIndex + len(crlf), requestLine, fmt.Errorf("Malformed request-line: %q\n", requestLineText)
	}

	method, target, version := parts[0], parts[1], parts[2]

	if !validMethod.MatchString(method) {
		return crlfIndex + len(crlf), requestLine, fmt.Errorf("Malformed method in request-line: %q\n", requestLineText)
	}

	if validHTTPVersion != version {
		return crlfIndex + len(crlf), requestLine, fmt.Errorf("Malformed HTTP version in request-line: %q\n", requestLineText)
	}
	versionParts := strings.Split(version, "/")
	httpVersion := versionParts[1]

	requestLine.Method = method
	requestLine.RequestTarget = target
	requestLine.HttpVersion = httpVersion

	return crlfIndex + len(crlf), requestLine, nil
}

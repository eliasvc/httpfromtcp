package request

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{}
	inputRequest, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Error reading request: %q\n", err.Error())
	}

	request.RequestLine, err = parseRequestLine(inputRequest)
	if err != nil {
		return nil, fmt.Errorf("Error parsing request-line: %q\n", err)
	}

	return &request, nil
}

func parseRequestLine(inputRequest []byte) (RequestLine, error) {
	var requestLine RequestLine
	validMethod := regexp.MustCompile(`[A-Z]+`)
	validHTTPVersion := "HTTP/1.1"

	inputRequestText := string(inputRequest)
	crlfIndex := strings.Index(inputRequestText, "\r\n")
	if crlfIndex == -1 {
		return requestLine, fmt.Errorf("CRLF not found on request-line")
	}
	requestLineText := inputRequestText[:crlfIndex]

	parts := strings.Split(requestLineText, " ")
	// request-line is composed of "method SP request-target SP HTTP-version"
	// where SP is single space. So it has to have only 3 parts if devided by
	// space
	if len(parts) != 3 {
		return requestLine, fmt.Errorf("Malformed request-line: %q\n", requestLineText)
	}

	method, target, version := parts[0], parts[1], parts[2]

	if !validMethod.MatchString(method) {
		return requestLine, fmt.Errorf("Malformed method in request-line: %q\n", requestLineText)
	}

	if validHTTPVersion != version {
		return requestLine, fmt.Errorf("Malformed HTTP version in request-line: %q\n", requestLineText)
	}
	versionParts := strings.Split(version, "/")
	httpVersion := versionParts[1]

	requestLine.Method = method
	requestLine.RequestTarget = target
	requestLine.HttpVersion = httpVersion

	return requestLine, nil
}

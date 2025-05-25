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
	sRequest := string(inputRequest)
	requestParts := strings.Split(sRequest, "\r\n")

	// '\r\n' not found in request
	if len(requestParts) == 1 {
		return nil, fmt.Errorf("Malformed request. '\\r\\n' not found: %q\n", sRequest)
	}

	request.RequestLine, err = parseRequestLine(requestParts[0])
	if err != nil {
		return nil, fmt.Errorf("Error parsing request-line: %q\n", err)
	}

	return &request, nil
}

func parseRequestLine(inputRequestLine string) (RequestLine, error) {
	var requestLine RequestLine
	validMethod := regexp.MustCompile(`[A-Z]+`)
	validHTTPVersion := "HTTP/1.1"

	parts := strings.Split(inputRequestLine, " ")
	// request-line is composed of "method SP request-target SP HTTP-version"
	// where SP is single space. So it has to have only 3 parts if devided by
	// space
	if len(parts) != 3 {
		return requestLine, fmt.Errorf("Malformed request-line: %q\n", inputRequestLine)
	}

	method, target, version := parts[0], parts[1], parts[2]

	if !validMethod.MatchString(method) {
		return requestLine, fmt.Errorf("Malformed method in request-line: %q\n", inputRequestLine)
	}

	requestLine.Method = method

	requestLine.RequestTarget = target

	if validHTTPVersion != version {
		return requestLine, fmt.Errorf("Malformed HTTP version in request-line: %q\n", inputRequestLine)
	}
	versionParts := strings.Split(version, "/")
	httpVersion := versionParts[1]
	requestLine.HttpVersion = httpVersion

	return requestLine, nil
}

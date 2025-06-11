package headers

import (
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	validHeader := regexp.MustCompile("^\\s*[A-Za-z0-9!#$%&'*+-.^_`|~]+$")
	sData := string(data)
	crlfIndex := strings.Index(sData, crlf)
	// CRLF not found. Assume we need more data
	if crlfIndex == -1 {
		return 0, false, nil
	}

	// Found CRLF that's between headers and body. We're done.
	if crlfIndex == 0 {
		return 0, true, nil
	}

	headerName, headerValue, found := strings.Cut(sData[:crlfIndex], ":")
	//fmt.Printf(">>header-name: %q, header-value: %q\n", headerName, headerValue)
	if !found {
		return 0, false, fmt.Errorf("Malformed header: %q\n", headerName)
	}

	if !validHeader.Match([]byte(headerName)) {
		return 0, false, fmt.Errorf("Malformed header: %q\n", headerName)
	}

	headerName = strings.TrimSpace(headerName)
	headerName = strings.ToLower(headerName)

	headerValue = strings.TrimSpace(headerValue)

	//fmt.Printf(">> Trimmed header-name: %q, header-value: %q\n", headerName, headerValue)
	// Combine field-values if the same field-name is already present
	if _, ok := h[headerName]; ok {
		h[headerName] = fmt.Sprintf("%s, %s", h[headerName], headerValue)
	} else {
		h[headerName] = headerValue
	}

	return crlfIndex + len(crlf), false, nil
}

func NewHeaders() Headers {
	return make(map[string]string)
}

package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
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
	if !found {
		return 0, false, fmt.Errorf("Malformed header: %q\n", headerName)
	}

	if strings.HasSuffix(headerName, " ") {
		return 0, false, fmt.Errorf("Malformed header: %q\n", headerName)
	}

	headerName = strings.TrimSpace(headerName)
	headerValue = strings.TrimSpace(headerValue)
	h[headerName] = headerValue

	return crlfIndex + len(crlf), false, nil
}

func NewHeaders() Headers {
	return Headers{}
}

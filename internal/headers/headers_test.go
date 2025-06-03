package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with valid extra whitespace in header name
	headers = NewHeaders()
	data = []byte("        Host: localhost:42069\r\n\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, 31, n)

	// Test: Valid 2 headers with existing headers
	headers = Headers{
		"Host":           "localhost",
		"Content-Length": "599",
	}
	data = []byte("Host: localhost:42069\r\nContent-Length: 69\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, 23, n)
	assert.False(t, done)
	assert.Equal(t, headers["Host"], "localhost:42069")
	assert.Equal(t, headers["Content-Length"], "599")

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, 0, n)
	assert.True(t, done)
	assert.NoError(t, err)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid field-name special characters
	headers = NewHeaders()
	data = []byte("!#$%&'*+-.^_`|~: boom\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid special characters
	headers = NewHeaders()
	data = []byte("!#$%&'*+-.^_`|~?><hahaha: boom\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

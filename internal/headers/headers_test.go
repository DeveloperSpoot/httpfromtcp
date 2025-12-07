package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeadersParse(t *testing.T) {

	// Test: Valid single header
	headers := make(Headers)
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = make(Headers)
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: invalid character in field name
	headers = make(Headers)
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid Character In Fildname with whitesapce.
	headers = make(Headers)
	data = []byte("         H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

package gelf

import (
    "testing"
)

type TSWriter struct {}
func (TSWriter) Write([]byte) (int, error) { return 0, nil }

func TestNewStream(t *testing.T) {
	_ = NewStream(&TSWriter{}, 0)		// Delimiter is "null char"
	_ = NewStream(&TSWriter{}, '\n')	// Delimiter is newline
}


package gelf

import (
	"io"
)

type stream struct {}

/*
	NewStream takes an io.Writer that writes to a stream connection
	(TCP, TLS, Unix domain socket) or anything similar (file, buf,
	whatever) and produces a GELFWriter.
*/
func NewStream(w io.Writer, delimiter byte) stream {
	return stream{}
}

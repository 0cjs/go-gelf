package gelf

import (
	"errors"
	"io"
)

type stream struct {
	w io.Writer
}

func newStream(w io.Writer, delimiter byte) stream {
	return stream{newDelimitedWriter(w, delimiter)}
}

func (s stream) Write(b []byte) (int, error) {
	return 0, nil
}

type delimitedWriter struct {
	w         io.Writer
	delimiter byte
}

func newDelimitedWriter(w io.Writer, d byte) delimitedWriter {
	return delimitedWriter{w, d}
}

func (dw delimitedWriter) Write(b []byte) (int, error) {
	n, err := dw.w.Write(b)
	if err != nil {
		return n, err
	}
	if len(b) != n {
		panic(errors.New("delimitedWriter: incomplete write"))
	}
	if n2, err := dw.w.Write([]byte{dw.delimiter}); err != nil {
		return n, err
	} else if n2 != 1 {
		panic(errors.New("delimitedWriter: unable to write delimiter"))
	}
	return n, nil
}

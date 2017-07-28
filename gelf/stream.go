package gelf

import (
	"io"

	"github.com/pkg/errors"
)

type stream struct {
	w io.Writer
}

func NewStream(w io.Writer, delimiter byte) stream {
	return stream{newDelimitedWriter(w, delimiter)}
}
func (s *stream) WriteMessage(m Message) error {
	_, err := s.w.Write(messageToJSON(m))
	return err
}

func (s *stream) Write(b []byte) (int, error) {
	if err := s.WriteMessage(messageFromByteSlice(b)); err != nil {
		return 0, err
	}
	return len(b), nil
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

package gelf

import "io"

type stream struct {
}

func NewStream(w io.Writer, delimiter byte) stream {
	return stream{}
}
func (s *stream) WriteMessage(m Message) error {
	return nil
}

func (s *stream) Write(b []byte) (int, error) {
	if err := s.WriteMessage(messageFromByteSlice(b)); err != nil {
		return 0, err
	}
	return len(b), nil
}

package gelf

import "io"

type stream struct {
	w io.Writer
}

func NewStream(w io.Writer, delimiter byte) stream {
	return stream{w}
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

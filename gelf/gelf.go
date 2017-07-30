package gelf

import (
	"io"
)

type GELFWriter interface {
	io.Writer
	WriteMessage(Message) error
}

type writer struct {
	w io.Writer
}

func NewStream(w io.Writer, delimiter byte) *writer {
	return &writer{newStream(w, delimiter)}
}

func NewPacket(w io.Writer, mtu uint32, compression Compression, level int) (*writer, error) {
	p, err := newPacket(w, mtu, compression, level)
	if err != nil {
		return nil, err
	}
	return &writer{p}, nil
}

func (w *writer) WriteMessage(m Message) error {
	_, err := w.w.Write(m.Bytes())
	return err
}

func (w *writer) Write(b []byte) (int, error) {
	if err := w.WriteMessage(messageFromByteSlice(b)); err != nil {
		return 0, err
	}
	return len(b), nil
}

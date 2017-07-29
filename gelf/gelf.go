package gelf

import (
	"errors"
	"io"
)

type GELFWriter interface {
	io.Writer
	WriteMessage(Message) error
}

type encoderType uint8

const (
	_ encoderType = iota
	streamEncoder
	packetEncoder
)

type encoderOptions struct {
	eType encoderType
	// Stream options
	delimiter byte
	// Packet options
	mtu         uint32
	compression Compression
	level       int
}

func StreamEncoder(d byte) encoderOptions {
	return encoderOptions{
		eType:     streamEncoder,
		delimiter: d,
	}
}

func PacketEncoder(mtu uint32, compression Compression, level int) encoderOptions {
	return encoderOptions{
		eType:       packetEncoder,
		mtu:         mtu,
		compression: compression,
		level:       level,
	}
}

type writer struct {
	w io.Writer
}

func NewWriter(w io.Writer, eo encoderOptions) (*writer, error) {
	switch eo.eType {
	case streamEncoder:
		return &writer{newStream(w, eo.delimiter)}, nil
	case packetEncoder:
		p, err := newPacket(w, eo.mtu, eo.compression, eo.level)
		if err != nil {
			return nil, err
		}
		return &writer{p}, nil
	}
	return nil, errors.New("invalid encoder type")
}

func (w *writer) WriteMessage(m Message) error {
	_, err := w.w.Write(messageToJSON(m))
	return err
}

func (w *writer) Write(b []byte) (int, error) {
	if err := w.WriteMessage(messageFromByteSlice(b)); err != nil {
		return 0, err
	}
	return len(b), nil
}

package gelf

import (
	"errors"
	"io"
)

type packet struct {
	w io.Writer
}

type Compression uint8

const (
	None Compression = iota
	Gzip
	Zlib
)

func NewPacket(w io.Writer, mtu uint32, compression Compression, level int) (*packet, error) {
	switch compression {
	case None:
		return &packet{newChunker(w)}, nil
	case Gzip:
		return newGzipPacketWriter(newChunker(w), level)
	case Zlib:
		return newZlibPacketWriter(newChunker(w), level)
	default:
		return nil, errors.New("invalid compression type")
	}

}

func (p *packet) WriteMessage(m Message) error {
	_, err := p.w.Write(messageToJSON(m))
	return err
}

func (p *packet) Write(b []byte) (int, error) {
	if err := p.WriteMessage(messageFromByteSlice(b)); err != nil {
		return 0, err
	}
	return len(b), nil
}

func newChunker(w io.Writer) io.Writer {
	return w
}

func newGzipPacketWriter(w io.Writer, level int) (*packet, error) {
	return &packet{w}, nil
}

func newZlibPacketWriter(w io.Writer, level int) (*packet, error) {
	return &packet{w}, nil
}

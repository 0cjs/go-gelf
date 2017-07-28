package gelf

import (
	"compress/gzip"
	"compress/zlib"
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
	chunker := newChunker(w, mtu)
	switch compression {
	case None:
		return &packet{chunker}, nil
	case Gzip:
		zw, err := gzip.NewWriterLevel(chunker, level)
		if err != nil {
			return nil, err
		}
		return &packet{newWriterSequence(chunker, zw)}, nil
	case Zlib:
		zw, err := zlib.NewWriterLevel(chunker, level)
		if err != nil {
			return nil, err
		}
		return &packet{newWriterSequence(chunker, zw)}, nil
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

func newChunker(w io.Writer, mtu uint32) io.Writer {
	return w
}

type writerSequence struct {
	streamedWriter  io.Writer
	restartedWriter WriteResetter
}

type WriteResetter interface {
	io.Writer
	Reset(io.Writer)
}

func newWriterSequence(streamedWriter io.Writer, restartedWriter WriteResetter) io.Writer {
	return writerSequence{streamedWriter, restartedWriter}
}

func (ws writerSequence) Write(b []byte) (int, error) {
	return 0, nil
}

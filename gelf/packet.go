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

func newPacket(w io.Writer, mtu uint32, compression Compression, level int) (*packet, error) {
	chunker := newChunker(w, mtu)
	switch compression {
	case None:
		return &packet{chunker}, nil
	case Gzip:
		zw, err := gzip.NewWriterLevel(nil, level)
		if err != nil {
			return nil, err
		}
		return &packet{newWriterSequence(chunker, zw)}, nil
	case Zlib:
		zw, err := zlib.NewWriterLevel(nil, level)
		if err != nil {
			return nil, err
		}
		return &packet{newWriterSequence(chunker, zw)}, nil
	default:
		return nil, errors.New("invalid compression type")
	}

}

func (p packet) Write(b []byte) (int, error) {
	return 0, nil
}

func newChunker(w io.Writer, mtu uint32) io.Writer {
	return w
}

type WriteCloseResetter interface {
	io.WriteCloser
	Reset(io.Writer)
}

type writerSequence struct {
	streamedWriter  io.Writer
	restartedWriter WriteCloseResetter
}

func newWriterSequence(streamedWriter io.Writer, restartedWriter WriteCloseResetter) io.Writer {
	return writerSequence{streamedWriter, restartedWriter}
}

func (ws writerSequence) Write(b []byte) (int, error) {
	ws.restartedWriter.Reset(ws.streamedWriter)
	defer ws.restartedWriter.Close()
	return ws.restartedWriter.Write(b)
}

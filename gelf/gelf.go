package gelf

import "io"

type stream struct{}

func NewStream(w io.Writer, delimiter byte) stream {
	return stream{}
}

type packet struct {
}

type Compression uint8

const (
	None Compression = iota
	Gzip
	Zlib
)

func NewPacket(w io.Writer, mtu uint32, compression Compression, level int) packet {
	return packet{}
}

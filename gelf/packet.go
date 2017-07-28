package gelf

import "io"

type packet struct {
	w io.Writer
}

type Compression uint8

const (
	None Compression = iota
	Gzip
	Zlib
)

func NewPacket(w io.Writer, mtu uint32, compression Compression, level int) packet {
	return packet{w}
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

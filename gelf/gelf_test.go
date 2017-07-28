package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"testing"
)

var buf bytes.Buffer

func TestNewStream(t *testing.T) {
	_ = NewStream(&buf, 0)
	_ = NewStream(&buf, '\n')
}

func TestNewPacket(t *testing.T) {
	_ = NewPacket(&buf, 1400, None, 0)
	_ = NewPacket(&buf, 1400, Gzip, gzip.BestSpeed)
	_ = NewPacket(&buf, 1400, Zlib, zlib.BestCompression)
	_ = NewPacket(&buf, 1234, 5, 0)
}

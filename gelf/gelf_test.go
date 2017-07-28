package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"testing"

	"github.com/stretchr/testify/assert"
)

var buf bytes.Buffer

func TestNewStream(t *testing.T) {
	_ = NewStream(&buf, 0)
	_ = NewStream(&buf, '\n')
}

func TestNewPacket(t *testing.T) {
	_ = NewPacket(&buf, 0, None, 0) // Let the library choose the MTU
	_ = NewPacket(&buf, 1400, Gzip, gzip.BestSpeed)
	_ = NewPacket(&buf, 1400, Zlib, zlib.BestCompression)
	_ = NewPacket(&buf, 1234, 5, 0)
}

func TestStream_Write(t *testing.T) {
	buf := new(bytes.Buffer)
	g := NewStream(buf, 0)
	data := "qwrtyuio"
	length, err := g.Write([]byte(data))
	assert.Nil(t, err)
	assert.Equal(t, len(data), length)
}

func TestPacket_Write(t *testing.T) {
	buf := new(bytes.Buffer)
	g := NewPacket(buf, 0, None, 0)
	data := "qwrtyuio"
	length, err := g.Write([]byte(data))
	assert.Nil(t, err)
	assert.Equal(t, len(data), length)
}

func TestMessageFromByteSlice(t *testing.T) {

}

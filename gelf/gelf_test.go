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
	_, err := NewPacket(&buf, 0, None, 0) // Let the library choose the MTU
	assert.Nil(t, err)
	_, err = NewPacket(&buf, 1400, Gzip, gzip.BestSpeed)
	assert.Nil(t, err)
	_, err = NewPacket(&buf, 1400, Zlib, zlib.BestCompression)
	assert.Nil(t, err)
	_, err = NewPacket(&buf, 1234, 5, 0)
	assert.EqualError(t, err, "invalid compression type")
	_, err = NewPacket(&buf, 1400, Gzip, -3)
	assert.EqualError(t, err, "gzip: invalid compression level: -3")
	_, err = NewPacket(&buf, 1400, Zlib, 10)
	assert.EqualError(t, err, "zlib: invalid compression level: 10")
}

func TestWriter_Write(t *testing.T) {
	buf := new(bytes.Buffer)
	g := NewStream(buf, 0)
	data := "qwrtyuio"
	length, err := g.Write([]byte(data))
	assert.Nil(t, err)
	assert.Equal(t, len(data), length)
}

func TestWriter_Write2(t *testing.T) {
	buf := new(bytes.Buffer)
	g, err := NewPacket(buf, 0, None, 0)
	assert.Nil(t, err)
	data := "qwrtyuio"
	length, err := g.Write([]byte(data))
	assert.Nil(t, err)
	assert.Equal(t, len(data), length)
}

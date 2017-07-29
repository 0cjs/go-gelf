package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var buf bytes.Buffer

func TestNewWriter(t *testing.T) {
	_, err := NewWriter(&buf, StreamEncoder(0))
	assert.Nil(t, err)
	_, err = NewWriter(&buf, StreamEncoder('\n'))
	assert.Nil(t, err)

	_, err = NewWriter(&buf, PacketEncoder(0, None, 0)) // Let the library choose the MTU
	assert.Nil(t, err)
	_, err = NewWriter(&buf, PacketEncoder(1400, Gzip, gzip.BestSpeed))
	assert.Nil(t, err)
	_, err = NewWriter(&buf, PacketEncoder(1400, Zlib, zlib.BestCompression))
	assert.Nil(t, err)
	_, err = NewWriter(&buf, PacketEncoder(1234, 5, 0))
	assert.EqualError(t, err, "invalid compression type")
	_, err = NewWriter(&buf, PacketEncoder(1400, Gzip, -3))
	assert.EqualError(t, err, "gzip: invalid compression level: -3")
	_, err = NewWriter(&buf, PacketEncoder(1400, Zlib, 10))
	assert.EqualError(t, err, "zlib: invalid compression level: 10")

	_, err = NewWriter(&buf, encoderOptions{})
	assert.EqualError(t, err, "invalid encoder type")
}

func TestWriter_Write(t *testing.T) {
	buf := new(bytes.Buffer)
	g, err := NewWriter(buf, StreamEncoder(0))
	assert.Nil(t, err)
	data := "qwrtyuio"
	length, err := g.Write([]byte(data))
	assert.Nil(t, err)
	assert.Equal(t, len(data), length)
}

func TestWriter_Write2(t *testing.T) {
	buf := new(bytes.Buffer)
	g, err := NewWriter(buf, PacketEncoder(0, None, 0))
	assert.Nil(t, err)
	data := "qwrtyuio"
	length, err := g.Write([]byte(data))
	assert.Nil(t, err)
	assert.Equal(t, len(data), length)
}

func TestDelimitedWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	dw := newDelimitedWriter(buf, '-')

	n, err := dw.Write([]byte("ab"))
	assert.Nil(t, err)
	assert.EqualValues(t, 2, n)

	n, err = dw.Write([]byte("c"))
	assert.Nil(t, err)
	assert.EqualValues(t, 1, n)

	n, err = dw.Write([]byte("def"))
	assert.Nil(t, err)
	assert.EqualValues(t, 3, n)

	assert.Equal(t, "ab-c-def-", buf.String())
}

type failingWriter struct{}

func (w failingWriter) Write(b []byte) (int, error) {
	return 1337, errors.New("failingWriter")
}
func TestDelimitedWriter_Write(t *testing.T) {
	dw := newDelimitedWriter(failingWriter{}, 'x')
	n, err := dw.Write([]byte("ab"))
	assert.EqualError(t, err, "failingWriter")
	assert.Equal(t, n, 1337)
}

type mockWriteResetter struct {
	w io.Writer
}

func (m mockWriteResetter) Write(b []byte) (int, error) {
	return m.w.Write(b)
}
func (m mockWriteResetter) Close() error {
	m.w.Write([]byte("RESET"))
	return nil
}

func TestWriterSequence_Write(t *testing.T) {
	buf := new(bytes.Buffer)
	ws := newWriterSequence(mockWriteResetter{buf})
	ws.Write([]byte("ab"))
	ws.Write([]byte("cd"))
	ws.Write([]byte("ef"))
	assert.Equal(t, "abRESETcdRESETefRESET", buf.String())
}

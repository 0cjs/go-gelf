package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockWriteResetter struct {
	w io.Writer
}

func (m mockWriteResetter) Write(b []byte) (int, error) {
	return m.w.Write(b)
}
func (m mockWriteResetter) Close() error {
	m.Write([]byte("---"))
	return nil
}

func (m mockWriteResetter) Reset(io.Writer) {
	m.w.Write([]byte("RESET"))
}

func TestWriterSequence_Write(t *testing.T) {
	buf := new(bytes.Buffer)
	ws := newWriterSequence(nil, mockWriteResetter{buf})
	ws.Write([]byte("ab"))
	ws.Write([]byte("cd"))
	ws.Write([]byte("ef"))
	assert.Equal(t, "RESETab---RESETcd---RESETef---", buf.String())
}

func TestWriterSequence_Write2(t *testing.T) {
	expected := "\x1f\x8b\b\x00\x00\x00\x00\x00\x04\xff\x00\x02\x00\xfd\xffab\x01\x00\x00\xff\xffmH\x83\x9e\x02\x00\x00\x00\x1f\x8b\b\x00\x00\x00\x00\x00\x04\xff\x00\x02\x00\xfd\xffcd\x01\x00\x00\xff\xffڏ\xd6E\x02\x00\x00\x00\x1f\x8b\b\x00\x00\x00\x00\x00\x04\xff\x00\x02\x00\xfd\xffef\x01\x00\x00\xff\xffpI\x82\xfd\x02\x00\x00\x00"

	buf := new(bytes.Buffer)
	zw, err := gzip.NewWriterLevel(nil, gzip.BestSpeed)
	assert.Nil(t, err)
	ws := newWriterSequence(buf, zw)
	ws.Write([]byte("ab"))
	ws.Write([]byte("cd"))
	ws.Write([]byte("ef"))
	assert.Equal(t, expected, buf.String())
}

func TestWriterSequence_Write3(t *testing.T) {
	expected := "x\x01\x00\x02\x00\xfd\xffab\x01\x00\x00\xff\xff\x01&\x00\xc4x\x01\x00\x02\x00\xfd\xffcd\x01\x00\x00\xff\xff\x01,\x00\xc8x\x01\x00\x02\x00\xfd\xffef\x01\x00\x00\xff\xff\x012\x00\xcc"

	buf := new(bytes.Buffer)
	zw, err := zlib.NewWriterLevel(nil, zlib.BestSpeed)
	assert.Nil(t, err)
	ws := newWriterSequence(buf, zw)
	ws.Write([]byte("ab"))
	ws.Write([]byte("cd"))
	ws.Write([]byte("ef"))
	assert.Equal(t, expected, buf.String())
}

package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"testing"

	"encoding/json"

	"github.com/pkg/errors"
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

func TestMessageToJSON(t *testing.T) {
	m := Message{
		Host:         "test-host",
		ShortMessage: "short message",
		Timestamp:    1234567890,
		Extra: map[string]interface{}{
			"foo": "bar",
		},
	}
	jsonBytes := messageToJSON(m)

	expected, _ := json.Marshal(map[string]interface{}{
		"host":          m.Host,
		"short_message": m.ShortMessage,
		"timestamp":     m.Timestamp,
		//"_foo":          "bar",
	})

	assert.JSONEq(t, string(expected), string(jsonBytes))
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

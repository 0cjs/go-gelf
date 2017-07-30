package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

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

func TestWriter_WriteMessage(t *testing.T) {
	buf := new(bytes.Buffer)
	w := NewStream(buf, 0)

	m := Message{
		Host:         "test-host",
		ShortMessage: "short message",
		FullMessage:  strings.Repeat("qwertyuiop", 500),
		Timestamp:    1234567890,
		Level:        6,
		Extra: map[string]interface{}{
			"foo": "bar",
		},
	}
	assert.Nil(t, w.WriteMessage(m))

	m2 := Message(m)
	m2.Timestamp = 1234567890
	m2.FullMessage = strings.Repeat("asdfgh", 50)
	assert.Nil(t, w.WriteMessage(m2))

	expected := bytes.NewBuffer(make([]byte, 0, len(m.Bytes())+len(m2.Bytes())+2))
	expected.Write(m.Bytes())
	expected.WriteByte(0)
	expected.Write(m2.Bytes())
	expected.WriteByte(0)
	assert.Equal(t, expected.String(), buf.String())
}

func BenchmarkStream_WriteMessage_WithExtra(b *testing.B) {
	hostname, err := os.Hostname()
	assert.Nil(b, err)

	w := NewStream(ioutil.Discard, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.WriteMessage(Message{
			Host:         hostname,
			ShortMessage: "short message",
			FullMessage:  "full message",
			Timestamp:    float64(time.Now().Unix()),
			Level:        6, // info
			Extra:        map[string]interface{}{"_file": "1234", "_line": "3456"},
		})
	}
}

func BenchmarkStream_WriteMessage_WithoutExtra(b *testing.B) {
	hostname, err := os.Hostname()
	assert.Nil(b, err)

	w := NewStream(ioutil.Discard, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.WriteMessage(Message{
			Host:         hostname,
			ShortMessage: "short message",
			FullMessage:  "full message",
			Timestamp:    float64(time.Now().Unix()),
			Level:        6, // info
		})
	}
}

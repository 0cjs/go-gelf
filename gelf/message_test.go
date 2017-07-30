package gelf

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"strings"

	"github.com/stretchr/testify/assert"
)

func TestMessageFromByteSlice(t *testing.T) {

}

func TestMessage_Bytes(t *testing.T) {
	m := Message{
		Host:         "test-host",
		ShortMessage: "short message",
		Timestamp:    1234567890,
		Extra: map[string]interface{}{
			"_foo": "bar",
		},
	}
	jsonBytes := m.Bytes()

	expected, _ := json.Marshal(map[string]interface{}{
		"host":          m.Host,
		"short_message": m.ShortMessage,
		"timestamp":     m.Timestamp,
		"_foo":          "bar",
	})

	assert.JSONEq(t, string(expected), string(jsonBytes))
}

func BenchmarkMessage_Bytes(b *testing.B) {
	host, err := os.Hostname()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := Message{
			Host:         host,
			ShortMessage: "short message",
			FullMessage:  "full message",
			Timestamp:    float64(time.Now().Unix()),
			Level:        6, // info
			Extra:        map[string]interface{}{"_file": "1234", "_line": "3456"},
		}
		m.Bytes()
	}
}

func BenchmarkMessage_Bytes2(b *testing.B) {
	host, err := os.Hostname()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := Message{
			Host:         host,
			ShortMessage: "short message",
			FullMessage:  "full message",
			Timestamp:    float64(time.Now().Unix()),
			Level:        6, // info
		}
		m.Bytes()
	}
}

func BenchmarkMessage_Bytes3(b *testing.B) {
	host, err := os.Hostname()
	if err != nil {
		b.Fatal(err)
	}
	longMessage := strings.Repeat("qwertyuiop", 500)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := Message{
			Host:         host,
			ShortMessage: "short message",
			FullMessage:  longMessage,
			Timestamp:    float64(time.Now().Unix()),
			Level:        6, // info
			Extra:        map[string]interface{}{"_file": "1234", "_line": "3456"},
		}
		m.Bytes()
	}
}

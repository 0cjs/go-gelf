package gelf

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

package gelf

import (
	"bytes"
	"encoding/json"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

func newBuffer() *bytes.Buffer {
	b := bufPool.Get().(*bytes.Buffer)
	if b != nil {
		b.Reset()
		return b
	}
	return bytes.NewBuffer(nil)
}

type Message struct {
	Host         string                 `json:"host"`
	ShortMessage string                 `json:"short_message"`
	FullMessage  string                 `json:"full_message,omitempty"`
	Timestamp    float64                `json:"timestamp"`
	Level        int32                  `json:"level,omitempty"`
	Extra        map[string]interface{} `json:"-"`
}

func messageFromByteSlice([]byte) Message {
	return Message{}
}

func (m Message) Bytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	hasExtra := len(m.Extra) > 0
	if !hasExtra {
		return b
	}

	eb, err := json.Marshal(m.Extra)
	if err != nil {
		panic(err)
	}

	buf := newBuffer()     // Get a buffer from the pool
	defer bufPool.Put(buf) // Put the buffer back into the pool when done

	buf.Write(b[:len(b)-1])      // Write the json of the message fields until the final closing bracket
	buf.WriteByte(',')           // Add a field separator
	buf.Write(eb[1 : len(eb)-1]) // Write the extra fields
	buf.WriteByte('}')           // Write the final closing bracket

	return buf.Bytes()
}

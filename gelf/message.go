package gelf

import "encoding/json"

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

func messageToJSON(m Message) []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return b
}

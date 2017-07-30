package gelf

import "os"

// Message represents the information that will be sent in a GELF
// payload to a logging server such as Graylog. All possible values
// are valid but some are not very sensible; use NewMessage() if you
// want reasonable defaults.
//
type Message struct {
	Host         string                 `json:"host"`
	ShortMessage string                 `json:"short_message"`
	FullMessage  string                 `json:"full_message,omitempty"`
	Timestamp    float64                `json:"timestamp"`
	Level        int32                  `json:"level,omitempty"`
	Extra        map[string]interface{} `json:"-"`
}

// Syslog severity levels valid for `Level` in a Message
// Specified by RFC 5424 https://tools.ietf.org/html/rfc5424
const (
	LOG_EMERG   = 0
	LOG_ALERT   = 1
	LOG_CRIT    = 2
	LOG_ERR     = 3
	LOG_WARNING = 4
	LOG_NOTICE  = 5
	LOG_INFO    = 6
	LOG_DEBUG   = 7
)

// NewMessage returns a Message with the fields filled in with
// appropriate defaults for the current host. These include the
// hostname of this host, the current timestamp, and a severity level
// of LOG_NOTICE.
// 
func NewMessage(shortMessage string) Message {
	m := Message{}
	// XXX It's not clear how to deal with failure here, or even
	//     what failure means.
	m.Host, _ = os.Hostname()
	m.ShortMessage = shortMessage
	m.Level = LOG_NOTICE
	return m
}

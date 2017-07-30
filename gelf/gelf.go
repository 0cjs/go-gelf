/*

Package `gelf` formats and writes Graylog GELF log messages to stream
or packet connections supplied by the caller. Information about GELF
can be found at http://docs.graylog.org/en/2.2/pages/gelf.html .

Though this logging protocol is definitively different from the syslog
protocol, the background on that from RFC 5424
(https://tools.ietf.org/html/rfc5424) can be useful to understand the
hows and whys of GELF.

*/
package gelf

import "io"

// GELFWriters are used to  send GELF messages to a server.
// Reception is not guranteed or confirmed.
//
// WriteMessage gives you full control of the message; NewMessage()
// will give you reasonable defaults to start with. If you're not too
// concerned with all the details of logging, io.Writer will take a
// string and fill in the rest for you.
//
type GELFWriter interface {
	io.Writer
	WriteMessage(Message) error
}

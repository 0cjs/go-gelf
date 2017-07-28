package gelf

type GELFWriter interface {
	Write([]byte) (int, error)
	WriteMessage(*Message) error
}

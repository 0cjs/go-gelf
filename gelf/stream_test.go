package gelf

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

type configurableWriter struct {
	fn func(b []byte) (int, error)
}

func (w *configurableWriter) Write(b []byte) (int, error) {
	return w.fn(b)
}

func TestDelimitedWriter_Write(t *testing.T) {
	failingWriter := configurableWriter{func(b []byte) (int, error) {
		return 1337, errors.New("failingWriter")
	}}
	dw := newDelimitedWriter(&failingWriter, 'x')
	n, err := dw.Write([]byte("ab"))
	assert.EqualError(t, err, "failingWriter")
	assert.Equal(t, 1337, n)
}

func TestDelimitedWriter_Write2(t *testing.T) {
	incompleteWriter := configurableWriter{func(b []byte) (int, error) {
		return len(b) / 2, nil
	}}
	dw := newDelimitedWriter(&incompleteWriter, 'x')
	assert.Panics(t, func() {
		dw.Write([]byte("ab"))
	})
}

func TestDelimitedWriter_Write3(t *testing.T) {
	failingDelimiterWriter := configurableWriter{func(b []byte) (int, error) {
		if len(b) == 1 {
			return 0, errors.New("failingDelimiterWriter")
		}
		return len(b), nil
	}}
	dw := newDelimitedWriter(&failingDelimiterWriter, 'x')
	n, err := dw.Write([]byte("ab"))
	assert.EqualError(t, err, "failingDelimiterWriter")
	assert.Equal(t, 2, n)
}

func TestDelimitedWriter_Write4(t *testing.T) {
	incompleteWriter := configurableWriter{func(b []byte) (int, error) {
		if len(b) == 1 {
			return 0, nil
		}
		return len(b), nil
	}}
	dw := newDelimitedWriter(&incompleteWriter, 'x')
	assert.Panics(t, func() {
		dw.Write([]byte("ab"))
	})
}

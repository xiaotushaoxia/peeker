package rwpeeker

import (
	"bytes"
	"io"
)

type PeekBytesWriter interface {
	PeekBytes() []byte
}

type Writer struct {
	pk  *bytes.Buffer
	ori io.Writer
	mw  io.Writer
}

func NewWriter(w io.Writer) *Writer {
	var w2 = bytes.NewBuffer([]byte{})
	writer := io.MultiWriter(w, w2)
	return &Writer{
		ori: w,
		pk:  w2,
		mw:  writer,
	}
}

func (w *Writer) Write(ps []byte) (int, error) {
	return w.mw.Write(ps)
}

func (w *Writer) WriteString(s string) (int, error) {
	return io.WriteString(w.mw, s)
}

func (w *Writer) PeekBytes() []byte {
	return w.pk.Bytes()
}

package rwpeeker

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

func PeekHttpRequestBody(c *gin.Context) ([]byte, error) {
	newBody, result := NewReader(c.Request.Body)
	c.Request.Body = newBody
	return result.All, result.Err
}

type GinResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *GinResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *GinResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (w *GinResponseWriter) PeekBytes() []byte {
	return w.body.Bytes()
}

func GinResponsePeeker(c *gin.Context) PeekBytesWriter {
	writer, ok := c.Writer.(PeekBytesWriter)
	if !ok {
		blw := &GinResponseWriter{body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
		c.Writer = blw
		writer = blw
	}
	return writer
}

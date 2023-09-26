package rwpeeker

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewWriter(t *testing.T) {
	var a = bytes.NewBuffer(nil)
	writer := NewWriter(a)

	writer.WriteString("1111")
	fmt.Println(writer.PeekBytes())
	writer.WriteString("2222")
	fmt.Println(writer.PeekBytes())
	writer.WriteString("3333")
	fmt.Println(writer.PeekBytes())
	writer.WriteString("4444")
	fmt.Println(writer.PeekBytes())

}

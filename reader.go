package rwpeeker

import (
	"bytes"
	"fmt"
	"io"
)

var ErrClosed = fmt.Errorf("use of closed reader")

type Reader struct {
	Body       *bytes.Reader
	Err        error
	CloseError error
	Total      []byte

	closed bool
}

// Read 尽量保持和原始的c.Request.Body的Read有一样的行为
// 因为用了io.ReadAll去读c.Request.Body， 所以可能是  Content = {"a": 1, "b  err = 连接被关闭
// 这时候r也能一样读到{"a": 1, "b，并返回在下一次再读的时候返回 err=连接被关闭
func (r *Reader) Read(p []byte) (int, error) {
	if r.closed {
		return 0, ErrClosed
	}
	n, err := r.Body.Read(p)
	if err != nil {
		if err != io.EOF { // bytes.Reader.Read不会返回除了io.EOF以外的错误
			panic(err) // 所以走到这里出问题了 panic吧
		}
		if r.Err == nil {
			return n, io.EOF
		}
		return n, r.Err
	}
	return n, nil
}
func (r *Reader) Close() error {
	r.closed = true
	return r.CloseError
}

func (r *Reader) Reset() {
	// 读完以后重置一下 可以给下一个人读
	r.Body = bytes.NewReader(r.Total)
}

type ReadAllResult struct {
	All []byte
	Err error
}

func NewReader(rc io.ReadCloser) (*Reader, *ReadAllResult) {
	requestBody, err := io.ReadAll(rc)

	result := &ReadAllResult{
		All: requestBody,
		Err: err,
	}
	//if err == nil {  // 测试 模拟超时用
	//	ctx, cc := context.WithTimeout(context.Background(), time.Millisecond)
	//	time.Sleep(time.Second / 100)
	//	err = ctx.Err()
	//	cc()
	//}

	// /net/http/request.go
	// For server requests, the Request Body is always non-nil
	// but will return EOF immediately when nobody is present.
	// The Server will close the request body. The ServeHTTP
	// Handler does not need to.
	// 其实是不需要close的  但是为了拿到close的err,这里close一下。相当于提前close了。没有关系
	return &Reader{
		Body:       bytes.NewReader(requestBody),
		Err:        err,
		Total:      requestBody,
		CloseError: rc.Close(),
	}, result
}

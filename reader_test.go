package rwpeeker

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestPeekReader(t *testing.T) {
	output1 := testPeekReader(&testReader{bs: []byte{1, 2, 3, 4, 5}}, false)
	output2 := testPeekReader(&testReader{bs: []byte{1, 2, 3, 4, 5}}, true)
	if !reflect.DeepEqual(output1, output2) {
		t.Errorf("NewReader is not equal to testReader: %v, %v", output1, output2)
	}

}

func TestTCPConnRead(t *testing.T) {
	// 测试net.Conn Read的行为
	listen, err3 := net.Listen("tcp", "127.0.0.1:8572")
	if err3 != nil {
		panic(err3)
	}

	var towrite []byte

	for i := 0; i < int(time.Now().Unix())%10+10; i++ {
		towrite = append(towrite, byte(i))
	}

	go func() {
		for i := 0; i < 2; i++ {
			accept, err := listen.Accept()
			if err != nil {
				panic(err)
			}
			nn, ee := accept.Write(towrite)
			fmt.Println("Accept one", nn, ee)
			accept.Close()
		}

		listen.Close()
		fmt.Println("close listen")
	}()

	dial1, err2 := net.Dial("tcp", "127.0.0.1:8572")
	if err2 != nil {
		panic(err2)
	}
	output1 := testPeekReader(dial1, false)

	dial2, err2 := net.Dial("tcp", "127.0.0.1:8572")
	if err2 != nil {
		panic(err2)
	}
	output2 := testPeekReader(dial2, true)

	if !reflect.DeepEqual(output1, output2) {
		t.Errorf("NewReader is not equal to testReader: %v, %v", output1, output2)
	}
}

func testPeekReader(rc io.ReadCloser, withPeek bool) string {
	var m = bytes.NewBuffer([]byte{})
	//var rc io.ReadCloser = &testReader{bs: []byte{1, 2, 3, 4, 5}}

	if withPeek {
		rc, _ = NewReader(rc)
	}

	var bs = make([]byte, 3)
	n, err := rc.Read(bs) // 3,nil
	fmt.Fprintf(m, "%v %v %v", bs, n, err)
	var bs2 = make([]byte, 10)
	n, err = rc.Read(bs2) // 2,nil  可以看到这里还是可以正常读出2字节，而不是返回  2,io.EOF
	fmt.Fprintf(m, "%v %v %v", bs, n, err)

	var bs3 = make([]byte, 3)
	n, err = rc.Read(bs3) // read EOF
	fmt.Fprintf(m, "%v %v %v", bs, n, err)

	return m.String()
}

type testReader struct {
	bs []byte
	i  int
}

func (r *testReader) Read(bs []byte) (int, error) {
	if r.i >= len(r.bs) {
		return 0, fmt.Errorf("read error")
	}
	copy(bs, r.bs[r.i:])
	cc := min(len(bs), len(r.bs[r.i:]))
	r.i += cc
	return cc, nil
}

func (r *testReader) Close() error {
	return nil
}
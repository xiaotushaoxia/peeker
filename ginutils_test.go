package rwpeeker

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestPeekHttpRequestBody(t *testing.T) {
	g := gin.New()
	g.Use(LoggerRW)

	g.POST("/xx", func(c *gin.Context) {
		var mp = map[string]any{}
		err := c.ShouldBindJSON(&mp)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}

		c.JSON(200, mp)
	})

	g.GET("/xx2", func(c *gin.Context) {
		c.JSON(200, "ok")
	})
	go func() {
		fmt.Println("exit", g.Run(":8484"))
	}()

	time.Sleep(time.Second)

	func() {
		resp, err := http.Get("http://127.0.0.1:8484/xx2")
		// Logger中间件打印
		// body:   err:  <nil>
		// resp:  "ok"
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			fmt.Println("get error", err)
			return
		}
		all, err := io.ReadAll(resp.Body)
		fmt.Println("resp", err, string(all))
	}()

	func() {
		resp, err := http.Post("http://127.0.0.1:8484/xx", "", bytes.NewBufferString("{\"a\":1}"))
		// Logger中间件打印
		// body:  {"a":1} err:  <nil>
		// resp:  {"a":1}
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			fmt.Println("get error", err)
			return
		}
		all, err := io.ReadAll(resp.Body)
		fmt.Println("resp", err, string(all))
	}()
}

func LoggerRW(c *gin.Context) {
	body, err := PeekHttpRequestBody(c)
	fmt.Println("body: ", string(body), "err: ", err)
	respPeeker := GinResponsePeeker(c)

	c.Next()

	fmt.Println("resp: ", string(respPeeker.PeekBytes()))
}

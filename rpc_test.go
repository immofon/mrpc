package mrpc

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/immofon/appoint/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TestRpc(t *testing.T) {
	log.TextMode()
	rpc := New(upgrader)

	rpc.RegisterFunc("echo", func(ctx context.Context, req Request) Return {
		msg := req.Get("msg", "please set field: msg")

		return req.Ret(Ok).Set("msg", msg)
	})

	http.Handle("/ws", rpc)
	log.L().Info("serve :8100")
	go func() {
		t.Error(http.ListenAndServe("localhost:8100", nil))
	}()

	time.Sleep(time.Millisecond * 100)

	fmt.Println("ok")
	c := NewClient("ws://localhost:8100/ws")
	go c.Serve()

	ret := c.Call(Req("echo").Set("msg", "hello mrpc"))
	fmt.Println(ret)
	ret = c.Call(Req("echo"))
	fmt.Println(ret)

}

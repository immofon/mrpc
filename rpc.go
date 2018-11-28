package mrpc

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/immofon/appoint/log"
)

type RPC struct {
	sync.Mutex
	upgrader websocket.Upgrader
	methods  map[string]Handler // key: method name
}

func New(upgrader websocket.Upgrader) *RPC {
	return &RPC{
		upgrader: upgrader,
		methods:  make(map[string]Handler, 0),
	}
}

func (rpc *RPC) Register(method string, handler Handler) {
	rpc.Lock()
	defer rpc.Unlock()
	if handler != nil {
		rpc.methods[method] = handler
	} else {
		delete(rpc.methods, method)
	}
}
func (rpc *RPC) RegisterFunc(method string, fn HandleFunc) {
	rpc.Register(method, fn)
}

func (rpc *RPC) Call(ctx context.Context, req Request) Return {
	rpc.Lock()
	defer rpc.Unlock()

	handler, ok := rpc.methods[req.Method]
	if !ok {
		log.L().WithField("name", req.Method).Warn("rpc method was not defined")
		return req.Ret(NotFound)
	}

	return handler.RPCHandle(ctx, req)
}

func (rpc *RPC) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := rpc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.E(err).Error()
		return
	}
	defer conn.Close()

	log.L().
		WithField("ip", conn.RemoteAddr().String()).
		Info("open connection")
	defer func() {
		log.L().
			WithField("ip", conn.RemoteAddr().String()).
			Info("close connection")
	}()

	ctx := context.Background()
	for {
		var req Request
		if err := conn.ReadJSON(&req); err != nil {
			log.E(err).Error()
			return
		}

		log.L().WithField("req", req).Info("get request")

		ret := rpc.Call(ctx, req)
		if ret.UpdateContext != nil {
			ctx = ret.UpdateContext(ctx)
		}

		if err := conn.WriteJSON(ret); err != nil {
			log.E(err).Error()
			return
		}
		log.L().
			WithField("ip", conn.RemoteAddr().String()).
			WithField("ret", ret).
			Debug("send response")

	}
}

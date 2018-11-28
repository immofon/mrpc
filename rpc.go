package rpc

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/immofon/appoint/log"
)

type Request struct {
	Id   string            `json:"id"`
	Func string            `json:"func"`
	Argv map[string]string `json:"argv"`
}

func (req Request) Get(key string, defaultv string) string {
	v, ok := req.Argv[key]
	if !ok {
		return defaultv
	}
	return v
}

func (req Request) Ret(status RetStatus) Return {
	return Return{
		Id:      req.Id,
		Status:  status,
		Details: make(map[string]string),
	}

}

type Return struct {
	Id            string                                `json:"id"`
	Status        RetStatus                             `json:"status"`
	Details       map[string]string                     `json:"details"`
	UpdateContext func(context.Context) context.Context `json:"-"`
}

func (ret Return) SetUpdateContext(fn func(context.Context) context.Context) Return {
	ret.UpdateContext = fn
	return ret
}
func (ret Return) Set(key, value string) Return {
	ret.Details[key] = value
	return ret
}
func (ret Return) Get(key, defaultv string) string {
	v, ok := ret.Details[key]
	if !ok {
		return defaultv
	}
	return v
}
func (ret Return) Has(key string) bool {
	_, ok := ret.Details[key]
	return ok
}

type Handler interface {
	RPCHandle(ctx context.Context, req Request) Return
}

type HandleFunc func(ctx context.Context, req Request) Return

func (hdf HandleFunc) RPCHandle(ctx context.Context, req Request) Return {
	return hdf(ctx, req)
}

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

	handler, ok := rpc.methods[req.Func]
	if !ok {
		log.L().WithField("name", req.Func).Warn("rpc method was not defined")
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

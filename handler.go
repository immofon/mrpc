package mrpc

import "context"

type Handler interface {
	RPCHandle(ctx context.Context, req Request) Return
}

type HandleFunc func(ctx context.Context, req Request) Return

func (hdf HandleFunc) RPCHandle(ctx context.Context, req Request) Return {
	return hdf(ctx, req)
}

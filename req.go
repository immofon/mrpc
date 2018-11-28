package mrpc

import (
	"context"
	"strings"
)

type Request struct {
	Id     string            `json:"id"`
	Method string            `json:"method"`
	Args   map[string]string `json:"args"`
}

func Req(method string) Request {
	return Request{
		Id:     "",
		Method: method,
		Args:   make(map[string]string),
	}
}

var encoder = strings.NewReplacer("π", "π0", ":", "π1")
var decoder = strings.NewReplacer("π1", ":", "π0", "π")

func (req Request) Get(key string, defaultv string) string {
	v, ok := req.Args[key]
	if !ok {
		return defaultv
	}
	return v
}

func (req Request) GetArray(key string, defaultvs []string) []string {
	raw := req.Get(key, "")
	if raw == "" {
		return defaultvs
	}

	vs := strings.Split(raw, ":")
	for i, v := range vs {
		vs[i] = decoder.Replace(v)
	}

	return vs
}

func (req Request) GetL(key string, defaultvs ...string) []string {
	return req.GetArray(key, defaultvs)
}

func (req Request) Set(k, v string) Request {
	req.Args[k] = v
	return req
}

func (req Request) SetArray(key string, vs []string) Request {
	for i, v := range vs {
		vs[i] = encoder.Replace(v)
	}
	return req.Set(key, strings.Join(vs, ":"))
}

func (req Request) SetL(key string, vs ...string) Request {
	return req.SetArray(key, vs)
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

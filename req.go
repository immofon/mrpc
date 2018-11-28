package mrpc

import "context"

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

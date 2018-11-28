package mrpc

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Request(t *testing.T) {
	req := Req("echo").SetL("something", "hello", "world", "π0π1π2ππ", "23π")
	jsonp(req)
	jsonp(req.GetL("something"))
}

func jsonp(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

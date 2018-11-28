package mrpc

import (
	"strconv"

	"github.com/gorilla/websocket"
)

type Callback func(Return)

type updateClientFn func(*Client)

type Client struct {
	ch chan updateClientFn

	conn      *websocket.Conn
	nextId    int
	callbacks map[string]Callback // key: id
}

func NewClient(url string) *Client {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(err)
		return nil
	}
	return &Client{
		ch: make(chan updateClientFn),

		conn:      conn,
		nextId:    1,
		callbacks: make(map[string]Callback),
	}
}

func (c *Client) Serve() {
	go func() {
		defer close(c.ch)
		defer c.conn.Close()

		var ret Return
		for {
			err := c.conn.ReadJSON(&ret)
			if err != nil {
				return
			}

			c.ch <- func(c *Client) {
				cb, ok := c.callbacks[ret.Id]
				if !ok {
					return
				}

				cb(ret)
				//TODO: subscribe don't delete callback function
				delete(c.callbacks, ret.Id)
			}
		}
	}()

	for fn := range c.ch {
		fn(c)
	}
}

func (c *Client) AsyncCall(req Request, cb Callback) {
	c.ch <- func(c *Client) {
		req.Id = c.generateId()
		c.callbacks[req.Id] = cb
		err := c.conn.WriteJSON(req)
		if err != nil {
			delete(c.callbacks, req.Id)
			cb(req.Ret(Network).Set("err", err.Error()))
		}
	}
}

func (c *Client) Call(req Request) Return {
	ch := make(chan Return)
	defer close(ch)

	c.AsyncCall(req, func(ret Return) {
		go func() {
			ch <- ret
		}()
	})
	return <-ch
}

func (c *Client) generateId() string {
	id := strconv.Itoa(c.nextId)
	c.nextId += 1
	return id
}

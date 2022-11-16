package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/silenceper/pool"
	"net"
	"time"
)

type Client struct {
	connPool pool.Pool
}

func NewClient(addr string) (*Client, error) {
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap: 10,
		MaxCap:     100,
		MaxIdle:    50,
		Factory: func() (interface{}, error) {
			return net.Dial("tcp", addr)
		},
		IdleTimeout: time.Minute,
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
	})
	if err != nil {
		return nil, err
	}
	return &Client{connPool: p}, nil
}

func (c Client) Invoke(ctx context.Context, req *Request) (*Response, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	obj, err := c.connPool.Get()
	if err != nil {
		return nil, err
	}
	conn := obj.(net.Conn)
	data = EncodeMsg(data)
	i, err := conn.Write(data)
	if i != len(data) {
		return nil, errors.New("micro: 未写入全部数据")
	}
	respMsg, err := ReadMsg(conn)
	if err != nil {
		return nil, err
	}
	return &Response{Data: respMsg}, nil
}

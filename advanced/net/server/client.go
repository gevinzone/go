package server

import (
	"encoding/binary"
	"net"
	"time"
)

type Client struct {
	addr string
}

func (c *Client) Send(msg string) (string, error) {
	conn, err := net.DialTimeout("tcp", c.addr, 3*time.Second)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = conn.Close()
	}()

	bs := make([]byte, lenBytes, len(msg)+lenBytes)
	binary.BigEndian.PutUint64(bs, uint64(len(msg)))
	bs = append(bs, msg...)

	_, err = conn.Write(bs)
	if err != nil {
		return "", err
	}

	lenBs := make([]byte, lenBytes)
	_, err = conn.Read(lenBs)
	if err != nil {
		return "", err
	}
	respLen := binary.BigEndian.Uint64(lenBs)
	resp := make([]byte, respLen)
	_, err = conn.Read(resp)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

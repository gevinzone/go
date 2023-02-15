package server2

import (
	"net"
	"time"
)

func Send(addr, msg string) (string, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*3)
	if err != nil {
		return "", err
	}
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return "", err
	}
	bs := make([]byte, 8)
	_, err = conn.Read(bs)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

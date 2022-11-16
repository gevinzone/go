package rpc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

const lenBytes = 8

func ReadMsg(conn net.Conn) (bs []byte, err error) {
	msgLenBytes := make([]byte, lenBytes)
	length, err := conn.Read(msgLenBytes)
	defer func() {
		if msg := recover(); msg != nil {
			err = errors.New(fmt.Sprintf("%v", msg))
		}
	}()
	if err != nil {
		return nil, err
	}
	if length != lenBytes {
		return nil, errors.New("could not read the length data")
	}
	dataLen := binary.BigEndian.Uint64(msgLenBytes)
	bs = make([]byte, dataLen)
	_, err = conn.Read(bs)
	return bs, err
}

func EncodeMsg(data []byte) []byte {
	l := len(data)
	resp := make([]byte, l+lenBytes)
	binary.BigEndian.PutUint64(resp, uint64(l))
	copy(resp, resp)
	return resp
}

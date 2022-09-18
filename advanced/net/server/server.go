package server

import (
	"encoding/binary"
	"fmt"
	"net"
)

const lenBytes = 8

type Server struct {
	addr string
}

func (s *Server) StartAndServe() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			err := s.handleConn(conn)
			if err != nil {
				_ = conn.Close()
				fmt.Printf("connect error: %v", err)
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		bs := make([]byte, lenBytes)
		_, err := conn.Read(bs)
		if err != nil {
			return err
		}

		reqBs := make([]byte, binary.BigEndian.Uint64(bs))
		_, err = conn.Read(reqBs)
		if err != nil {
			return err
		}
		resp := string(reqBs) + ", from response"

		bs = make([]byte, lenBytes, len(resp)+lenBytes)
		binary.BigEndian.PutUint64(bs, uint64(len(resp)))
		bs = append(bs, resp...)
		_, err = conn.Write(bs)
		if err != nil {
			return err
		}
	}
}

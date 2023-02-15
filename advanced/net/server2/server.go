package server2

import (
	"io"
	"net"
)

func StartAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			if err := handleConn(conn); err != nil {
				conn.Close()
			}
		}()
	}
}

func handleConn(conn net.Conn) error {
	bs := make([]byte, 8)
	_, err := conn.Read(bs)
	if err != nil {
		return err
	}
	_, err = conn.Write(bs)
	if err == net.ErrClosed || err == io.EOF || err == io.ErrUnexpectedEOF {
		return err
	}
	return nil
}

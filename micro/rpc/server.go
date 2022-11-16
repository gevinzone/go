package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

type Server struct {
	services map[string]reflectionStub
}

func NewServer() *Server {
	return &Server{services: map[string]reflectionStub{}}
}

func (s *Server) Register(service Service) error {
	s.services[service.Name()] = reflectionStub{
		value: reflect.ValueOf(service),
	}
	return nil
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			// log
			// return
			continue
		}
		go func() {
			if err := s.handleConn(conn); err != nil {
				_ = conn.Close()
				return
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		reqMsg, err := ReadMsg(conn)
		if err != nil {
			return err
		}
		req := &Request{}
		err = json.Unmarshal(reqMsg, req)
		if err != nil {
			return err
		}
		service, ok := s.services[req.ServiceName]
		if !ok {
			return errors.New("找不到服务")
		}
		ctx := context.Background()
		data, err := service.Invoke(ctx, req.MethodName, reqMsg)
		if err != nil {
			return err
		}
		data = EncodeMsg(data)
		_, err = conn.Write(data)
		return err
	}
}

type reflectionStub struct {
	value reflect.Value
}

func (r *reflectionStub) Invoke(ctx context.Context, methodName string, data []byte) ([]byte, error) {
	method := r.value.MethodByName(methodName)
	inType := method.Type().In(1)
	in := reflect.New(inType.Elem())
	err := json.Unmarshal(data, in.Interface())
	if err != nil {
		return nil, err
	}
	res := method.Call([]reflect.Value{reflect.ValueOf(ctx), in})
	if len(res) > 1 && !res[1].IsZero() {
		return nil, res[1].Interface().(error)
	}
	return json.Marshal(res[0].Interface())
}

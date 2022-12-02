package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/gevinzone/go/rpc/basic/proto/gen"
	"google.golang.org/grpc"
	"io"
	"net"
	"strconv"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "启动端口号")
	flag.Parse()
}

type GreeterServer struct {
}

func (s *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello, %s", r.Name)}, nil
}

func (s *GreeterServer) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
	for n := 0; n <= 6; n++ {
		_ = stream.Send(&pb.HelloReply{Message: "hello.list" + strconv.Itoa(n)})
		_ = stream.Send(&pb.HelloReply{Message: r.Name})
	}

	return nil

}
func (s *GreeterServer) SayRecord(stream pb.Greeter_SayRecordServer) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReply{Message: "say.record"})
		}
		if err != nil {
			return err
		}

		fmt.Printf("resp: %v\n", resp.Name)
	}

}
func (s *GreeterServer) SayRoute(stream pb.Greeter_SayRouteServer) error {
	n := 0
	for {
		_ = stream.Send(&pb.HelloReply{Message: "say.route"})
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n++
		fmt.Printf("resp: %v\n", resp)
	}
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	listener, _ := net.Listen("tcp", ":"+port)
	_ = server.Serve(listener)
}

package main

import (
	"flag"
	pb "github.com/gevinzone/go/grpc/tag-service/proto/gen"
	"github.com/gevinzone/go/grpc/tag-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8004", "启动端口号")
	flag.Parse()
}

func main() {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("server.Serve err: %v", err)
	}
}

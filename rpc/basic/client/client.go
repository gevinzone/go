package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/gevinzone/go/rpc/basic/proto/gen"
	"google.golang.org/grpc"
	"io"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "服务端口号")
}

func main() {
	conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	_ = SayHello(client)
	_ = SayList(client, &pb.HelloRequest{Name: "Gevin"})
	_ = SayRecord(client, &pb.HelloRequest{Name: "Gevin"})
	_ = SayRoute(client, &pb.HelloRequest{Name: "Gevin"})
}

func SayHello(client pb.GreeterClient) error {
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Gevin"})
	fmt.Printf("client.SayHello resp: %s", resp.Message)
	return nil
}

func SayList(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, _ := client.SayList(context.Background(), r)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("resp: %v\n", resp.Message)
	}
}

func SayRecord(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, _ := client.SayRecord(context.Background())
	for n := 0; n < 6; n++ {
		_ = stream.Send(r)
	}
	resp, _ := stream.CloseAndRecv()
	fmt.Printf("resp err: %v\n", resp)
	return nil
}

func SayRoute(client pb.GreeterClient, r *pb.HelloRequest) error {
	//stream, _ := client.SayRoute(context.Background())
	//for n := 0; n <= 6; n++ {
	//	_ = stream.Send(r)
	//	resp, err := stream.Recv()
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		return err
	//	}
	//
	//	log.Printf("resp err: %v", resp)
	//}
	//
	//_ = stream.CloseSend()
	//
	//return nil
	stream, _ := client.SayRoute(context.Background())
	for i := 0; i < 6; i++ {
		_ = stream.Send(r)
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("resp: %v\n", resp)
	}
	_ = stream.CloseSend()
	return nil
}

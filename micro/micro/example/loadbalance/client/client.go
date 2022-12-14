package main

import (
	"context"
	"fmt"
	"github.com/gevinzone/go/micro/micro"
	"github.com/gevinzone/go/micro/micro/example/loadbalance/proto/gen"
	"github.com/gevinzone/go/micro/micro/loadbalancer"
	"github.com/gevinzone/go/micro/micro/loadbalancer/roundrobin"
	"github.com/gevinzone/go/micro/micro/registry/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"log"
	"time"
)

func main() {
	// 注册中心
	// 注册中心
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		panic(err)
	}
	r, err := etcd.NewRegistry(etcdClient)
	if err != nil {
		panic(err)
	}
	// 注册你的负载均衡策略
	pickerBuilder := &roundrobin.PickerBuilder{
		Filter: loadbalancer.GroupFilter,
	}
	builder := base.NewBalancerBuilder(pickerBuilder.Name(), pickerBuilder, base.Config{HealthCheck: true})
	balancer.Register(builder)

	cc, err := grpc.Dial("registry:///user-service",
		grpc.WithInsecure(),
		grpc.WithResolvers(micro.NewResolverBuilder(r, time.Second*3)),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`,
			pickerBuilder.Name())))
	if err != nil {
		panic(err)
	}
	client := gen.NewUserServiceClient(cc)
	for i := 0; i < 100; i++ {
		ctx := context.WithValue(context.Background(), "group", "b")
		resp, err := client.GetById(ctx, &gen.GetByIdReq{})
		if err != nil {
			panic(err)
		}
		log.Println(resp.User)
	}
}

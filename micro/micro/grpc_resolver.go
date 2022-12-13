package micro

import (
	"context"
	"github.com/gevinzone/go/micro/micro/registry"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

type grpcResolverBuilder struct {
	r       registry.Registry
	timeout time.Duration
}

func NewResolverBuilder(r registry.Registry, timeout time.Duration) resolver.Builder {
	return &grpcResolverBuilder{
		r:       r,
		timeout: timeout,
	}
}

func (g *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	res := &grpcResolver{
		target:  target,
		cc:      cc,
		r:       g.r,
		close:   make(chan struct{}, 1),
		timeout: g.timeout,
	}
	state := res.resolve()
	log.Println(state)
	return res, res.watch()
}

func (g *grpcResolverBuilder) Scheme() string {
	return "r"
}

type grpcResolver struct {
	target  resolver.Target
	cc      resolver.ClientConn
	r       registry.Registry
	close   chan struct{}
	timeout time.Duration
}

func (g *grpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
	g.resolve()
}

func (g *grpcResolver) Close() {
	g.close <- struct{}{}
}

func (g *grpcResolver) watch() error {
	eventsCh, err := g.r.Subscribe(g.target.URL.Path)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case event := <-eventsCh:
				g.resolve()
				log.Println(event)
			case <-g.close:
				close(g.close)
				return
			}
		}
	}()
	return nil
}

func (g *grpcResolver) resolve() resolver.State {
	r := g.r
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	instances, err := r.ListService(ctx, g.target.URL.Path)
	if err != nil {
		g.cc.ReportError(err)
		return resolver.State{}
	}
	address := make([]resolver.Address, 0, len(instances))
	for _, ins := range instances {
		address = append(address, newAddress(ins))
	}
	state := resolver.State{
		Addresses: address,
	}
	err = g.cc.UpdateState(state)
	if err != nil {
		g.cc.ReportError(err)
	}
	return state
}

func newAddress(ins registry.ServiceInstance) resolver.Address {
	return resolver.Address{
		Addr:       ins.Address,
		ServerName: ins.ServiceName,
		Attributes: attributes.New("weight", ins.Weight).
			WithValue("group", ins.Group),
	}
}

package roundrobin

import (
	"github.com/gevinzone/go/micro/micro/loadbalancer"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"sync/atomic"
)

type Picker struct {
	ins    []instance
	cnt    uint64
	filter loadbalancer.Filter
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	candidates := make([]instance, 0, len(p.ins))
	for _, sub := range p.ins {
		if !p.filter(info, sub.address) {
			continue
		}
		candidates = append(candidates, sub)
	}
	if len(candidates) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	cnt := atomic.AddUint64(&p.cnt, 1)
	index := cnt % uint64(len(candidates))
	return balancer.PickResult{
		SubConn: candidates[index].sub,
		Done: func(info balancer.DoneInfo) {
			if info.Err != nil {

			}
		},
	}, nil
}

type PickerBuilder struct {
	Filter loadbalancer.Filter
}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]instance, 0, len(info.ReadySCs))
	for sub, subInfo := range info.ReadySCs {
		conns = append(conns, instance{
			sub:     sub,
			address: subInfo.Address,
		})
	}
	return &Picker{
		ins:    conns,
		filter: p.Filter,
	}
}

func (p *PickerBuilder) Name() string {
	return "ROUND_ROBIN"
}

type instance struct {
	sub     balancer.SubConn
	address resolver.Address
}

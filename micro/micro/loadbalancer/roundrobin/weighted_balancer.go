package roundrobin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
)

type WeightPicker struct {
	mutex sync.Mutex
	conns []*conn
}

func (w *WeightPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(w.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	w.mutex.Lock()
	defer w.mutex.Unlock()
	var (
		totalWeight   uint32 = 0
		maxWeightConn *conn
	)
	for _, c := range w.conns {
		efficientWeight := c.efficientWeight
		totalWeight += efficientWeight
		c.currentWeight = c.currentWeight + efficientWeight
		if maxWeightConn == nil || maxWeightConn.currentWeight < c.currentWeight {
			maxWeightConn = c
		}
		maxWeightConn.currentWeight = maxWeightConn.currentWeight - totalWeight
	}
	return balancer.PickResult{
		SubConn: maxWeightConn.SubConn,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type WeightBuilder struct {
}

func (w *WeightBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*conn, 0, len(info.ReadySCs))
	for subConn, subConninfo := range info.ReadySCs {
		weight := uint32(subConninfo.Address.Attributes.Value("weight").(int))
		conns = append(conns, &conn{
			SubConn:         subConn,
			weight:          weight,
			currentWeight:   weight,
			efficientWeight: weight,
		})
	}
	return &WeightPicker{
		conns: conns,
	}
}

func (*WeightBuilder) Name() string {
	return "WEIGHT_ROUND_ROBIN"
}

type conn struct {
	balancer.SubConn
	weight          uint32
	currentWeight   uint32
	efficientWeight uint32
}

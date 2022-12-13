package registry

import "context"

//go:generate mockgen -package=mocks -destination=mocks/registry.mock.go -source=types.go Registry
type Registry interface {
	Register(ctx context.Context, ins ServiceInstance) error
	Unregister(ctx context.Context, ins ServiceInstance) error
	ListService(ctx context.Context, serviceName string) ([]ServiceInstance, error)
	Subscribe(serviceName string) (<-chan Event, error)
	Close() error
}

type ServiceInstance struct {
	ServiceName string
	Address     string
	Weight      uint32
	Group       string
}

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeAdd
	EventTypeDelete
	EventTypeUpdate
)

type Event struct {
	Type     EventType
	Instance ServiceInstance
}

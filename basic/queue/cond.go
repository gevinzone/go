package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type Cond struct {
	signal chan struct{}
	l      sync.Locker
}

func NewCond(l sync.Locker) *Cond {
	return &Cond{
		signal: make(chan struct{}),
		l:      l,
	}
}

func (c *Cond) Broadcast() {
	var old chan struct{}
	c.signal, old = make(chan struct{}), c.signal
	c.l.Unlock()
	close(old)
}

func (c *Cond) SignalCh() <-chan struct{} {
	res := c.signal
	c.l.Unlock()
	return res
}

type Condition struct {
	L sync.Locker
	n unsafe.Pointer
}

func NewCondition(l sync.Locker) *Condition {
	ch := make(chan struct{})
	return &Condition{
		L: l,
		n: unsafe.Pointer(&ch),
	}
}

func (c *Condition) Wait() {
	ch := c.NotifyChan()
	c.L.Unlock()
	<-ch
	c.L.Lock()
}

func (c *Condition) WaitWithCtx(ctx context.Context) error {
	ch := c.NotifyChan()
	c.L.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		c.L.Lock()
		return nil
	}
}

func (c *Condition) WaitWithTimeout(t time.Duration) {
	n := c.NotifyChan()
	c.L.Unlock()
	select {
	case <-n:
	case <-time.After(t):
	}
	c.L.Lock()
}

func (c *Condition) Broadcast() {
	n := make(chan struct{})
	oldPtr := atomic.SwapPointer(&c.n, unsafe.Pointer(&n))
	close(*((*chan struct{})(oldPtr)))
}

func (c *Condition) NotifyChan() <-chan struct{} {
	p := atomic.LoadPointer(&c.n)
	return *((*chan struct{})(p))
}

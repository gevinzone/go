package queue

import "sync"

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

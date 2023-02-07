package syncx

import "github.com/google/uuid"

type ChLock struct {
	ch  chan string
	val string
}

func NewChLock() *ChLock {
	val := uuid.New().String()
	res := &ChLock{
		ch:  make(chan string, 1),
		val: val,
	}
	res.ch <- val
	return res
}

func (c *ChLock) Lock() string {
	<-c.ch
	return c.val
}

func (c *ChLock) Unlock(id string) {
	if id != c.val {
		panic("invalid not unlock")
	}
	val := uuid.New().String()
	select {
	case c.ch <- val:
		c.val = val
		return
	default:
		panic("not locked")
	}
}

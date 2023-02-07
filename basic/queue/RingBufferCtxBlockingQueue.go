package queue

import (
	"context"
	"sync"
)

type RingBufferCtxBlockingQueue[T any] struct {
	q                   *queue[T]
	m                   *sync.RWMutex
	queueNotEmptySignal *Cond
	queueNotFullSignal  *Cond
}

func NewRingBufferBlockingQueueV2[T any](capacity int) *RingBufferCtxBlockingQueue[T] {
	l := &sync.RWMutex{}
	return &RingBufferCtxBlockingQueue[T]{
		q:                   newQueue[T](capacity),
		m:                   l,
		queueNotEmptySignal: NewCond(l),
		queueNotFullSignal:  NewCond(l),
	}
}

func (r *RingBufferCtxBlockingQueue[T]) IsFull() bool {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.q.IsFull()
}

func (r *RingBufferCtxBlockingQueue[T]) IsEmpty() bool {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.q.IsEmpty()
}

func (r *RingBufferCtxBlockingQueue[T]) Enqueue(ctx context.Context, val T) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	r.m.Lock()
	for r.IsFull() {
		signal := r.queueNotFullSignal.SignalCh()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-signal:
			r.m.Lock()
		}
	}
	err := r.q.Enqueue(val)
	r.queueNotEmptySignal.Broadcast()
	return err
}

func (r *RingBufferCtxBlockingQueue[T]) Dequeue(ctx context.Context) (T, error) {
	var t T
	if ctx.Err() != nil {
		return t, ctx.Err()
	}
	r.m.Lock()
	for r.IsEmpty() {
		signal := r.queueNotEmptySignal.SignalCh()
		select {
		case <-ctx.Done():
			return t, ctx.Err()
		case <-signal:
			r.m.Lock()
		}
	}
	t, err := r.q.Dequeue()
	r.queueNotFullSignal.Broadcast()
	return t, err
}

func (r *RingBufferCtxBlockingQueue[T]) Peek() (T, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.q.Peek()
}

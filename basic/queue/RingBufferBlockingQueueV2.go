package queue

import (
	"sync"
)

type RingBufferBlockingQueueV2[T any] struct {
	q                   *queue[T]
	m                   *sync.RWMutex
	queueNotEmptySignal *sync.Cond
	queueNotFullSignal  *sync.Cond
}

func NewRingBufferBlockingQueueV2[T any](capacity int) *RingBufferBlockingQueueV2[T] {
	m := &sync.RWMutex{}
	return &RingBufferBlockingQueueV2[T]{
		q:                   newQueue[T](capacity),
		m:                   m,
		queueNotEmptySignal: sync.NewCond(m),
		queueNotFullSignal:  sync.NewCond(m),
	}
}

func (r *RingBufferBlockingQueueV2[T]) IsFull() bool {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.q.IsFull()
}

func (r *RingBufferBlockingQueueV2[T]) IsEmpty() bool {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.q.IsEmpty()
}

func (r *RingBufferBlockingQueueV2[T]) Enqueue(val T) error {
	r.m.Lock()
	for r.q.IsFull() {
		r.queueNotFullSignal.Wait()
	}
	err := r.q.Enqueue(val)
	r.queueNotEmptySignal.Broadcast()
	r.m.Unlock()
	return err
}

func (r *RingBufferBlockingQueueV2[T]) Dequeue() (T, error) {
	r.m.Lock()
	for r.q.IsEmpty() {
		r.queueNotEmptySignal.Wait()
	}
	val, err := r.q.Dequeue()
	r.queueNotFullSignal.Broadcast()
	r.m.Unlock()
	return val, err
}

func (r *RingBufferBlockingQueueV2[T]) Peek() (T, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.q.Peek()
}

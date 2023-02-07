package queue

import (
	"sync"
)

type RingBufferBlockingQueue[T any] struct {
	q                   *queue[T]
	m                   *sync.RWMutex
	queueNotEmptySignal chan struct{}
	queueNotFullSignal  chan struct{}
}

func NewRingBufferBlockingQueue[T any](capacity int) *RingBufferBlockingQueue[T] {
	return &RingBufferBlockingQueue[T]{
		q:                   newQueue[T](capacity),
		queueNotEmptySignal: make(chan struct{}, 1),
		queueNotFullSignal:  make(chan struct{}, 1),
	}
}

func (r *RingBufferBlockingQueue[T]) IsFull() bool {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.q.IsFull()
}

func (r *RingBufferBlockingQueue[T]) IsEmpty() bool {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.q.IsEmpty()
}

func (r *RingBufferBlockingQueue[T]) Enqueue(val T) error {
	r.m.Lock()
	for r.q.IsFull() {
		r.m.Unlock()
		<-r.queueNotFullSignal
		r.m.Lock()
	}
	err := r.q.Enqueue(val)
	if len(r.queueNotEmptySignal) == 0 {
		r.queueNotEmptySignal <- struct{}{}
	}
	r.m.Unlock()
	return err
}

func (r *RingBufferBlockingQueue[T]) Dequeue() (T, error) {
	r.m.Lock()
	for r.q.IsEmpty() {
		r.m.Unlock()
		<-r.queueNotEmptySignal
		r.m.Lock()
	}
	val, err := r.q.Dequeue()
	if len(r.queueNotFullSignal) == 0 {
		r.queueNotFullSignal <- struct{}{}
	}
	r.m.Unlock()
	return val, err
}

func (r *RingBufferBlockingQueue[T]) Peek() (T, error) {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.q.Peek()
}

package queue

import (
	"context"
	"reflect"
	"sync/atomic"
	"time"
	"unsafe"
)

type ConcurrentLinkedQueue[T any] struct {
	head  unsafe.Pointer
	tail  unsafe.Pointer
	count int64
}

func (c *ConcurrentLinkedQueue[T]) Enqueue(ctx context.Context, val T) error {
	n := node[T]{val: val}
	nPtr := unsafe.Pointer(&n)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			tail := atomic.LoadPointer(&c.tail)
			if atomic.CompareAndSwapPointer(&c.tail, tail, nPtr) {
				tn := *(*node[T])(tail)
				//tn.next = nPtr
				atomic.StorePointer(&tn.next, nPtr)
				atomic.AddInt64(&c.count, 1)
				return nil
			}
			time.Sleep(10 * time.Second)

		}
	}
}

func (c *ConcurrentLinkedQueue[T]) Dequeue(ctx context.Context) (T, error) {
	for {
		select {
		case <-ctx.Done():
			var t T
			return t, ctx.Err()
		default:
			head := atomic.LoadPointer(&c.head)
			hn := (*node[T])(head)
			tailPtr := atomic.LoadPointer(&c.tail)
			tn := (*node[T])(tailPtr)
			if hn == tn {
				var t T
				return t, errQueueEmpty
			}
			nextPtr := atomic.LoadPointer(&hn.next)
			if atomic.CompareAndSwapPointer(&c.head, head, nextPtr) {
				atomic.AddInt64(&c.count, -1)
				return hn.val, nil
			}
		}
	}
}

func (c *ConcurrentLinkedQueue[T]) Peek() (T, error) {
	var zero T
	count := atomic.LoadInt64(&c.count)
	if count == 0 {
		return zero, errQueueEmpty
	}
	headPtr := atomic.LoadPointer(&c.head)
	head := (*node[T])(headPtr)
	if reflect.ValueOf(head).IsZero() {
		return zero, errQueueEmpty
	}
	return head.val, nil
}

var _ BlockingQueue[int] = &ConcurrentLinkedQueue[int]{}

type node[T any] struct {
	next unsafe.Pointer
	val  T
}

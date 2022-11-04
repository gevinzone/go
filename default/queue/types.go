package queue

import "context"

type Queue[T any] interface {
	Enqueue(val T) error
	Dequeue() (T, error)
	Peek() (T, error)
}

type BlockingQueue[T any] interface {
	Enqueue(ctx context.Context, val T) error
	Dequeue(ctx context.Context) (T, error)
	Peek() (T, error)
}

package queue

import "errors"

var (
	errQueueFull  = errors.New("queue: queue is full")
	errQueueEmpty = errors.New("queue: queue is empty")
)

// queue 队列，其中data为队列的总容量，是一个ring buffer，head指向队头，tail指向下一个入队元素的位置
// 为了区分空队列和满队列的情况，留空一个位置，即head == tail 时，表示空队列，head == tail + 1时，表示满队列
// 队列实际可用容量为cap(data) - 1
type queue[T any] struct {
	data []T
	head int
	tail int
}

func newQueue[T any](capacity int) *queue[T] {
	return &queue[T]{
		data: make([]T, capacity+1),
		head: 0,
		tail: 0,
	}
}

func (q *queue[T]) IsFull() bool {
	return (q.tail+1)%cap(q.data) == q.head
}

func (q *queue[T]) IsEmpty() bool {
	return q.head == q.tail
}

func (q *queue[T]) cap() int {
	return cap(q.data) - 1
}

func (q *queue[T]) Enqueue(val T) error {
	if q.IsFull() {
		return errQueueFull
	}
	q.data[q.tail] = val
	q.tail++
	if q.tail == cap(q.data) {
		q.tail = 0
	}
	return nil
}

func (q *queue[T]) Dequeue() (T, error) {
	var t T
	if q.IsEmpty() {
		return t, errQueueEmpty
	}
	t, q.data[q.head] = q.data[q.head], t
	q.head++
	if q.head == cap(q.data) {
		q.head = 0
	}
	return t, nil
}

func (q *queue[T]) Peek() (T, error) {
	var t T
	if q.IsEmpty() {
		return t, errQueueEmpty
	}
	t = q.data[q.head]
	return t, nil
}

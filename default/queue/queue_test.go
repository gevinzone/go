package queue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewQueue(t *testing.T) {
	capacity := 10
	q := newQueue[int](capacity)
	assert.Equal(t, capacity, q.cap())
}

func TestQueue_Enqueue(t *testing.T) {
	testCases := []struct {
		name      string
		q         *queue[int]
		val       int
		wantHead  int
		wantTail  int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "empty queue",
			q:         createQueue(5),
			val:       1,
			wantHead:  0,
			wantTail:  1,
			wantSlice: []int{1, 0, 0, 0, 0, 0},
		},
		{
			name:    "full queue",
			q:       createQueue(3, 1, 2, 3),
			val:     1,
			wantErr: errQueueFull,
		},
		{
			name:      "normal queue",
			q:         createQueue(5, 1, 2, 3),
			val:       4,
			wantHead:  0,
			wantTail:  4,
			wantSlice: []int{1, 2, 3, 4, 0, 0},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.q.Enqueue(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantHead, tc.q.head)
			assert.Equal(t, tc.wantTail, tc.q.tail)
			assert.Equal(t, tc.wantSlice, tc.q.data)
		})
	}
}

func TestQueue_Dequeue(t *testing.T) {
	testCases := []struct {
		name      string
		q         *queue[int]
		val       int
		wantHead  int
		wantTail  int
		wantSlice []int
		wantErr   error
	}{
		{
			name:    "empty queue",
			q:       createQueue(5),
			wantErr: errQueueEmpty,
		},
		{
			name:      "full queue",
			q:         createQueue(3, 1, 2, 3),
			val:       1,
			wantHead:  1,
			wantTail:  3,
			wantSlice: []int{0, 2, 3, 0},
		},
		{
			name:      "normal queue",
			q:         createQueue(5, 1, 2, 3),
			val:       1,
			wantHead:  1,
			wantTail:  3,
			wantSlice: []int{0, 2, 3, 0, 0, 0},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.q.Dequeue()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.val, val)
			assert.Equal(t, tc.wantHead, tc.q.head)
			assert.Equal(t, tc.wantTail, tc.q.tail)
			assert.Equal(t, tc.wantSlice, tc.q.data)
		})
	}
}

func TestQueue_Peek(t *testing.T) {
	testCases := []struct {
		name    string
		q       *queue[int]
		val     int
		wantErr error
	}{
		{
			name:    "empty queue",
			q:       createQueue(5),
			wantErr: errQueueEmpty,
		},
		{
			name: "full queue",
			q:    createQueue(3, 1, 2, 3),
			val:  1,
		},
		{
			name: "normal queue",
			q:    createQueue(5, 1, 2, 3),
			val:  1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.q.Peek()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.val, val)
		})
	}
}

func createQueue(capacity int, vals ...int) *queue[int] {
	data := make([]int, capacity+1)
	for i, val := range vals {
		data[i] = val
	}
	q := &queue[int]{
		data: data,
		head: 0,
		tail: len(vals),
	}
	return q
}

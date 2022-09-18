package queue

import (
	"testing"
	"time"
)

func TestTaskQueue(t *testing.T) {
	tq := TaskQueue{
		ch:        make(chan int, 10),
		consumers: []TaskConsumer{},
	}

	t.Log(tq)
	t.Log(tq.ch)

	tq.RegisterConsumer(Consumer1)
	tq.RegisterConsumer(Consumer2)
	t.Log(tq.consumers)

	tq.Start()

	for i := 0; i < 100; i++ {
		tq.PublishTask(i)
	}
	tq.Close()
	time.Sleep(2 * time.Second)

}

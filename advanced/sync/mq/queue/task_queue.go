package queue

import (
	"fmt"
	"time"
)

type TaskConsumer func(taskChan *chan int)
type TaskQueue struct {
	ch        chan int
	consumers []TaskConsumer
}

func (tq *TaskQueue) RegisterConsumer(consumer TaskConsumer) {
	tq.consumers = append(tq.consumers, consumer)
}

func (tq *TaskQueue) PublishTask(i int) {
	tq.ch <- i
}

func (tq *TaskQueue) Start() {
	for _, consumer := range tq.consumers {
		consumer(&tq.ch)
	}
}

func (tq *TaskQueue) Close() {
	close(tq.ch)
}

func Consumer1(taskChan *chan int) {
	go func() {
		for {
			i, ok := <-*taskChan
			if !ok {
				return
			}
			time.Sleep(time.Millisecond)
			fmt.Printf("Consumer1: %d\n", i)
		}
	}()
}

func Consumer2(taskChan *chan int) {
	go func() {
		for {
			i, ok := <-*taskChan
			if !ok {
				return
			}
			time.Sleep(time.Millisecond)
			fmt.Printf("Consumer2: %d\n", i)
		}
	}()
}

package mq

import (
	"fmt"
	"testing"
	"time"
)

func TestProducerConsumer(t *testing.T) {

	broker := Broker{Consumers: make([]Consumer, 0)}
	limit := 2
	for i := 0; i < limit; i++ {
		c := CreateConsumer()
		broker.Register(c)
		c.Consume()
	}
	for i := 0; i < limit*5; i++ {
		broker.Publish(i)
	}

	time.Sleep(3 * time.Second)
	broker.Close()

}

func CreateConsumer() *Consumer1 {
	return &Consumer1{
		BaseConsumer{
			Ch: make(chan int, 10),
		},
	}
}

func TestBroker2(t *testing.T) {
	broker := NewBroker2()
	limit := 3
	for i := 0; i < limit; i++ {
		broker.Register(func(s string) {
			fmt.Println(fmt.Sprintf("msg: %s", s))
		})
	}
	broker.Start()
	for i := 0; i < limit*3; i++ {
		broker.produce(fmt.Sprintf("msg%d", i))
	}
	time.Sleep(3 * time.Second)
}

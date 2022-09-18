package mq

import "fmt"

type Broker struct {
	Consumers []Consumer
}

func (b *Broker) Register(c Consumer) {
	b.Consumers = append(b.Consumers, c)
}

func (b *Broker) Close() {
	for _, consumer := range b.Consumers {
		consumer.Stop()
	}
}

func (b *Broker) Publish(data int) {
	for _, consumer := range b.Consumers {
		consumer.Subscribe(data)
	}
}

type Consumer interface {
	Consume()
	Stop()
	Subscribe(data int)
}

type BaseConsumer struct {
	Ch chan int
}

func (c *BaseConsumer) Subscribe(data int) {
	c.Ch <- data
}

type Consumer1 struct {
	BaseConsumer
}

func (c *Consumer1) Consume() {
	go func() {
		for {
			data, ok := <-c.Ch
			if !ok {
				return
			}
			fmt.Println(data)
		}
	}()

}

func (c *Consumer1) Stop() {
	close(c.Ch)

}

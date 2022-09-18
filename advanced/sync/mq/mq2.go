package mq

type Broker2 struct {
	ch        chan string
	consumers []func(s string)
}

func (b *Broker2) produce(msg string) {
	b.ch <- msg
}

func (b *Broker2) Register(consumer func(s string)) {
	b.consumers = append(b.consumers, consumer)
}

func (b *Broker2) Start() {
	go func() {
		for {
			s, ok := <-b.ch
			if !ok {
				return
			}
			for _, c := range b.consumers {
				c(s)
			}
		}
	}()
}

func NewBroker2() *Broker2 {
	return &Broker2{
		ch:        make(chan string, 10),
		consumers: make([]func(s string), 0, 10),
	}
}

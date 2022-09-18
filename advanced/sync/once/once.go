package once

import (
	"fmt"
	"sync"
)

type user struct {
	Name string
}

type AdminInstance struct {
	once           sync.Once
	SingleInstance user
}

func (receiver *AdminInstance) getInstance() *user {
	receiver.once.Do(func() {
		receiver.SingleInstance = user{
			Name: "admin",
		}
	})
	return &receiver.SingleInstance
}

type CloseOnlyOnce struct {
	close sync.Once
}

func (o *CloseOnlyOnce) Close() error {
	o.close.Do(func() {
		fmt.Println("close")
	})
	return nil
}

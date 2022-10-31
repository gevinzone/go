package channel

import (
	"fmt"
	"time"
)

func orchestrate() {
	lockChan := make(chan struct{}, 1)
	res := 1
	lockChan <- struct{}{}
	for i := 0; i < 4; i++ {
		go func() {
			<-lockChan
			fmt.Println(res)
			res++
			//time.Sleep(time.Second)
			lockChan <- struct{}{}
		}()
	}
	time.Sleep(1 * time.Second)
}

func orchestrate2() {
	lockChan := make(chan struct{})
	res := 1
	for i := 0; i < 4; i++ {
		go func() {
			<-lockChan
			fmt.Println(res)
			res++
			//time.Sleep(time.Second)
			lockChan <- struct{}{}
		}()
	}
	lockChan <- struct{}{}
	time.Sleep(1 * time.Second)
}

func orchestrate3() {
	channels := [...]chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}
	for i := 0; i < 4; i++ {
		go func(i int) {
			<-channels[i]
			fmt.Println(i + 1)
			channels[(i+1)%4] <- struct{}{}
		}(i)
	}
	channels[0] <- struct{}{}
	select {
	case <-channels[0]:
		fmt.Println("complete")
	}
}
func orchestrate4() {
	channels := [...]chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}
	for i := 0; i < 4; i++ {
		go func(i int) {
			for j := 0; j < 10; j++ {
				token := <-channels[i]
				fmt.Println(i + 1)
				//time.Sleep(time.Second)
				channels[(i+1)%4] <- token
			}

		}(i)
	}
	channels[0] <- struct{}{}
	select {}
}

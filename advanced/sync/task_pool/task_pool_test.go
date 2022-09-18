package task_pool

import (
	"fmt"
	"testing"
	"time"
)

func TestTaskPoolWithToken_Do(t *testing.T) {
	limit := 2
	tp := NewTaskPoolWithToken(limit)
	for i := 0; i < limit; i++ {
		tp.Do(func() {
			time.Sleep(time.Second)
			fmt.Printf("complete task...\n")
		})
	}

	tp.Do(func() {
		fmt.Println("complete task")
	})
	time.Sleep(2 * time.Second)
}

func TestTaskPoolWithCache_Do(t *testing.T) {
	limit := 2
	tp := NewTaskPoolWithCache(limit)
	tp.Start()
	//tp.Stop()
	//tp.Restart()
	for i := 0; i < limit; i++ {
		tp.Do(func() {
			time.Sleep(time.Second)
			fmt.Printf("complete task...\n")
		})
	}

	tp.Do(func() {
		fmt.Println("complete task")
	})
	time.Sleep(2 * time.Second)
	tp.Stop()
}

package practice

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContextWithValue(t *testing.T) {
	bg := context.Background()
	levelOneCtx := context.WithValue(bg, "k1", "v1")
	levelTwoCtx := context.WithValue(levelOneCtx, "k2", "v2")

	t.Logf(levelOneCtx.Value("k1").(string))
	t.Logf("%t", levelOneCtx.Value("k2") == nil)
	t.Logf(levelTwoCtx.Value("k1").(string))
	t.Logf(levelTwoCtx.Value("k2").(string))
	err := levelTwoCtx.Err()
	if err != nil {
		t.Log(err)
	}
}

func TestContextTimeout(t *testing.T) {
	bg := context.Background()
	timeoutCtx, cancel := context.WithTimeout(bg, time.Second*1)
	defer cancel()
	subTimeoutCtx, subCancel := context.WithTimeout(timeoutCtx, time.Second*3)

	go func() {
		<-subTimeoutCtx.Done()
		t.Log("timeout")
		if err := subTimeoutCtx.Err(); err != nil {
			t.Logf("error: %v", err)
		}
	}()

	time.Sleep(time.Second * 2)
	subCancel()
}

func TestTimeoutExample(t *testing.T) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	bsChan := make(chan struct{})

	go func() {
		slowBusiness()
		bsChan <- struct{}{}
	}()
	select {
	case <-timeoutCtx.Done():
		t.Log("business run timeout")
	case <-bsChan:
		t.Log("business complete")
	}
}

func slowBusiness() {
	time.Sleep(time.Second * 2)
	fmt.Println("work done!")
}

func TestTimeoutWithTimeAfter(t *testing.T) {
	bsChan := make(chan struct{})
	go func() {
		slowBusiness()
		bsChan <- struct{}{}
	}()

	timer := time.AfterFunc(time.Second, func() {
		t.Log("business run timeout")
	})
	<-bsChan
	timer.Stop()
}

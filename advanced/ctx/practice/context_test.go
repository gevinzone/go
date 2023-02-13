package practice

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
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

func TestContextTimeoutSelect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	waitTime := time.Millisecond
	err := selective(ctx, waitTime)
	assert.NoError(t, err)

	// 进入default，虽然执行过程中已经超时，但依然会执行完，不会跳到另一个分支
	waitTime = time.Second * 2
	err = selective(ctx, waitTime)
	assert.NoError(t, err)

	// 这是ctx已经超时了
	waitTime = time.Millisecond
	err = selective(ctx, waitTime)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestContextTimeoutSelect2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// select 只要进入一个分支，就是执行到结束

	waitTime := time.Millisecond
	err := selective2(ctx, waitTime)
	assert.NoError(t, err)

	waitTime = time.Millisecond * 700
	err = selective2(ctx, waitTime)
	assert.NoError(t, err)

	// 这时ctx已经超时
	waitTime = time.Millisecond
	err = selective2(ctx, waitTime)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func selective(ctx context.Context, t time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		fmt.Println("processing...")
		time.Sleep(t)
		fmt.Println("done.")
		return nil
	}
}

func selective2(ctx context.Context, t time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(t):
		fmt.Println("processing...")
		time.Sleep(t * 2)
		fmt.Println("done.")
		return nil
	}
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

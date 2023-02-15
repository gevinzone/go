package server2

import (
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	addr := ":8080"
	go func() {
		err := StartAndServe(addr)
		if err != nil {
			t.Log(err)
		}
	}()
	time.Sleep(time.Second * 3)
	msg := "hello"
	resp, err := Send(addr, msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

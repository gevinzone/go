package server

import "testing"

func TestClient_Send(t *testing.T) {
	client := &Client{addr: "127.0.0.1:8080"}
	resp, err := client.Send("hello, world")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}

package server

import "testing"

func TestServer_StartAndServe(t *testing.T) {
	server := &Server{addr: ":8080"}
	_ = server.StartAndServe()
}

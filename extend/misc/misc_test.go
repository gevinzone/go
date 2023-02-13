package misc

import (
	"fmt"
	"testing"
)

type function struct {
	f func()
}

func TestFunction(t *testing.T) {
	f := &function{f: func() {
		fmt.Println("hello")
	}}
	callFunction(f)

	f = &function{}
	callFunction(f)
}

func callFunction(f *function) {
	if f.f == nil {
		fmt.Println("f.f is nil")
		return
	}
	f.f()
}

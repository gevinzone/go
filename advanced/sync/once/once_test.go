package once

import (
	"testing"
)

func TestSingleInstance(t *testing.T) {
	admin := AdminInstance{}
	for i := 0; i < 5; i++ {
		instance := admin.getInstance()
		t.Logf("instance address is %p", instance)
	}
}

func TestCloseOnce(t *testing.T) {
	once := CloseOnlyOnce{}
	for i := 0; i < 5; i++ {
		err := once.Close()
		if err != nil {
			t.Error(err)
		}
	}
}

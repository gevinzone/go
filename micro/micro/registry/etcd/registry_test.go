package etcd

import (
	"github.com/gevinzone/go/micro/micro/registry/etcd/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	clientV3 "go.etcd.io/etcd/client/v3"
	"log"
	"testing"
)

func TestRegistry_Subscribe(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func() clientV3.Watcher
		wantErr error
	}{
		{
			name: "normal",
			mock: func() clientV3.Watcher {
				watcher := mocks.NewMockWatcher(&gomock.Controller{})
				return watcher
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &Registry{
				client: &clientV3.Client{
					Watcher: tc.mock(),
				},
			}
			ch, err := r.Subscribe("service-name")
			assert.Equal(t, tc.wantErr, err)
			event := <-ch
			log.Println(event)
		})
	}
}

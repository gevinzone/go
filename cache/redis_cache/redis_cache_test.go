package redis_cache

import (
	"context"
	redis_mock "github.com/gevinzone/go/cache/mock"
	"github.com/go-redis/redis/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCache_Set(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	testCases := []struct {
		name    string
		mock    func() redis.Cmdable
		key     string
		val     any
		wantErr error
	}{
		{
			name: "success",
			key:  "key1",
			val:  123,
			mock: func() redis.Cmdable {
				client := redis_mock.NewMockCmdable(ctl)
				cmd := redis.NewStatusCmd(context.Background())
				cmd.SetVal("OK")
				client.EXPECT().
					Set(gomock.Any(), "key1", 123, time.Second).
					Return(cmd)
				return client
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := RedisCache{
				client: tc.mock(),
			}
			err := c.Set(context.Background(), tc.key, tc.val, time.Second)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

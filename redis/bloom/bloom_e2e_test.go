package bloom

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBloomExample(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	// install redis stack first
	assert.NoError(t, bloomExample(context.Background(), c))
}

package local_cache

import (
	"context"
	"errors"
	"github.com/gevinzone/go/cache"
	"sync"
	"sync/atomic"
	"time"
)

type MaxCntCache struct {
	mutex  sync.Mutex
	MaxCnt int32
	Cnt    int32
	*LocalCache
}

func NewMaxCntCache(maxCnt int32) *MaxCntCache {
	res := &MaxCntCache{MaxCnt: maxCnt}
	res.LocalCache = NewLocalCache(func(key string, val any) {
		atomic.AddInt32(&res.Cnt, -1)
	})
	return res
}

func (m *MaxCntCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, err := m.Get(ctx, key)
	if err != nil && err != cache.ErrKeyNotFound {
		return err
	}
	if err == cache.ErrKeyNotFound {
		cnt := atomic.AddInt32(&m.Cnt, 1)
		if cnt > m.MaxCnt {
			atomic.AddInt32(&m.Cnt, -1)
			return errors.New("cache: 已经满了")
		}
	}
	return m.LocalCache.Set(ctx, key, val, expiration)
}

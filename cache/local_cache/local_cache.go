package local_cache

import (
	"context"
	"github.com/gevinzone/go/cache/internal/errs"
	"sync"
	"time"
)

type LocalCache struct {
	data         map[string]any
	mutex        sync.RWMutex
	closeChan    chan struct{}
	closeOnce    sync.Once
	perClearLoop int

	onEvicted func(key string, val any)
}

func NewLocalCache(onEvicted func(key string, val any)) *LocalCache {
	res := &LocalCache{
		closeChan:    make(chan struct{}),
		perClearLoop: 1000,
	}
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				cnt := 0
				res.mutex.Lock()
				for k, v := range res.data {
					if v.(*item).deadline.Before(t) {
						res.delete(k, v)
					}
					cnt++
					if cnt >= res.perClearLoop {
						break
					}
				}
				res.mutex.Unlock()
			case <-res.closeChan:
				return
			}
		}
	}()
	return res
}

func (l *LocalCache) Get(ctx context.Context, key string) (any, error) {
	l.mutex.RLock()
	val, ok := l.data[key]
	if !ok {
		return nil, errs.NewErrKeyNotFound(key)
	}
	l.mutex.RUnlock()
	itm := val.(*item)
	// 懒/延迟删除
	if itm.deadline.Before(time.Now()) {
		l.mutex.Lock()
		defer l.mutex.Unlock()
		val, ok = l.data[key]
		if !ok {
			return nil, errs.NewErrKeyNotFound(key)
		}
		itm = val.(*item)
		if itm.deadline.Before(time.Now()) {
			l.delete(key, itm.val)
		}
		return nil, errs.NewErrKeyNotFound(key)
	}
	return itm.val, nil
}

// Set 用于为存储更新键值对，如果expiration=0，则不会过期
func (l *LocalCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.data[key] = &item{
		val:      val,
		deadline: time.Now().Add(expiration),
	}
	return nil
}

func (l *LocalCache) Delete(ctx context.Context, key string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	val, ok := l.data[key]
	if !ok {
		return nil
	}
	l.delete(key, val.(*item).val)
	return nil
}

func (l *LocalCache) delete(key string, val any) {
	delete(l.data, key)
	if l.onEvicted != nil {
		l.onEvicted(key, val)
	}

}

func (l *LocalCache) Close() error {
	l.closeOnce.Do(func() {
		l.closeChan <- struct{}{}
		close(l.closeChan)
	})
	return nil
}

type item struct {
	val      any
	deadline time.Time
}

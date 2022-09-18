package object_pool

import (
	"sync"
	"unsafe"
)

type ObjectPool interface {
	get() any
	put(val any)
}

type MyPool struct {
	pool     sync.Pool
	capacity int
	count    int
	maxSize  uintptr
}

func (m *MyPool) get() any {
	return m.pool.Get()
}

func (m *MyPool) put(val any) {
	if unsafe.Sizeof(val) > m.maxSize {
		return
	}

	m.pool.Put(val)
}


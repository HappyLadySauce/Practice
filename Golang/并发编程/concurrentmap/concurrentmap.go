package concurrentmap

import (
	"sync"
)

type ConcurrentMap[K comparable,V any] struct {
	data	map[K]V
	mu 		sync.RWMutex
}

func NewConcurrentMap[K comparable, V any](cap int) *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		data: make(map[K]V, cap),
	}
}

func (m *ConcurrentMap[K, V]) Store(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *ConcurrentMap[K, V]) Load(key K) (value V, exists bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists = m.data[key]
	return
}



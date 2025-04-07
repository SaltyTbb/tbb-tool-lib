package datastructures

import "sync"

// 这是一个阻塞读的map
type SyncMap[K comparable, V any] struct {
	mu       sync.RWMutex
	items    map[K]V
	condiMap map[K]*sync.Cond
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		mu:       sync.RWMutex{},
		items:    make(map[K]V),
		condiMap: make(map[K]*sync.Cond),
	}
}

func (m *SyncMap[K, V]) Get(key K) (V, bool) {
	m.mu.Lock()

	item, ok := m.items[key]
	if ok {
		m.mu.Unlock()
		return item, true
	}

	var cond *sync.Cond
	if c, exists := m.condiMap[key]; exists {
		cond = c
	} else {
		cond = sync.NewCond(&m.mu)
		m.condiMap[key] = cond
	}

	cond.Wait()

	item = m.items[key]
	m.mu.Unlock()
	return item, true
}

func (m *SyncMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = value

	if cond, ok := m.condiMap[key]; ok {
		cond.Broadcast()
		delete(m.condiMap, key)
	}
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, key)
}

func (m *SyncMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.items)
}

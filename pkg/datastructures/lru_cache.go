package datastructures

import (
	"container/list"
	"fmt"
)

type LRUCache[K comparable, V any] struct {
	cache    map[K]*cacheItem[V]
	capacity int
	ll       *list.List
}

type cacheItem[T any] struct {
	node  *list.Element
	value T
}

func NewLruCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		cache:    make(map[K]*cacheItem[V]),
		capacity: capacity,
		ll:       list.New(),
	}
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	item, ok := c.cache[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.ll.MoveToFront(item.node)
	return item.value, true
}

func (c *LRUCache[K, V]) Put(key K, value V) {
	item, ok := c.cache[key]
	if ok {
		c.ll.MoveToFront(item.node)
		item.value = value
		return
	}
	if c.ll.Len() == c.capacity {
		last := c.ll.Back()
		c.ll.Remove(last)
		delete(c.cache, last.Value.(K))
	}
	c.cache[key] = &cacheItem[V]{
		node:  c.ll.PushFront(key),
		value: value,
	}
	c.ll.PushFront(key)
}

func (c *LRUCache[K, V]) Remove(key K) {
	item, ok := c.cache[key]
	if !ok {
		return
	}
	c.ll.Remove(item.node)
	delete(c.cache, key)
}

func (c *LRUCache[K, V]) Len() int {
	return c.ll.Len()
}

func (c *LRUCache[K, V]) Clear() {
	c.ll.Init()
	c.cache = make(map[K]*cacheItem[V])
}

func (c *LRUCache[K, V]) String() string {
	return fmt.Sprintf("LRUCache{cache: %v, capacity: %d, len: %d}", c.cache, c.capacity, len(c.cache))
}

func (v *cacheItem[T]) String() string {
	return fmt.Sprintf("cacheItem{value: %v}", v.value)
}

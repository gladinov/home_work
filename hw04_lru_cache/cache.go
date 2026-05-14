package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value any
}

func newCacheItem(key Key, value any) *cacheItem {
	return &cacheItem{
		key:   key,
		value: value,
	}
}

func NewCache(capacity int) Cache {
	if capacity < 0 {
		capacity = 0
	}
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value any) bool {
	if c.capacity <= 0 {
		return false
	}
	cacheItemValue := newCacheItem(key, value)
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.items[key]
	if ok {
		elem.Value = cacheItemValue
		c.queue.MoveToFront(elem)
		return true
	}
	c.deleteBack()
	item := c.queue.PushFront(cacheItemValue)
	c.items[key] = item
	return false
}

func (c *lruCache) Get(key Key) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(elem)
	item, ok := elem.Value.(*cacheItem)
	if !ok {
		panic("lru cache invariant violation: list item value is not *cacheItem")
	}
	return item.value, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) deleteBack() {
	if c.capacity <= 0 || c.queue.Len() < c.capacity {
		return
	}
	back := c.queue.Back()
	cacheItemValue, ok := back.Value.(*cacheItem)
	if !ok {
		panic("lru cache invariant violation: list item value is not *cacheItem")
	}
	delete(c.items, cacheItemValue.key)
	c.queue.Remove(back)
}

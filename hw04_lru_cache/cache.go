package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
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
	elem, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(elem)
	item := elem.Value.(*cacheItem)
	return item.value, true
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) deleteBack() {
	if c.capacity <= 0 || c.queue.Len() < c.capacity {
		return
	}
	back := c.queue.Back()
	cacheItemValue := back.Value.(*cacheItem)

	delete(c.items, cacheItemValue.key)
	c.queue.Remove(back)
}

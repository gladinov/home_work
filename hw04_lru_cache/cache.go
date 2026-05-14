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
	value List
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value any) bool {
	elem, ok := c.items[key]
	if ok {
		elem.Value = value
		c.queue.MoveToFront(elem)
		return true
	}
	c.deleteBack()
	item := c.queue.PushFront(key, value)
	c.items[key] = item
	return false
}

func (c *lruCache) deleteBack() {
	if c.queue.Len() < c.capacity {
		return
	}
	delete(c.items, c.queue.Back().Key)
	c.queue.Remove(c.queue.Back())
}

func (c *lruCache) Get(key Key) (any, bool) {
	elem, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(elem)
		return elem.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

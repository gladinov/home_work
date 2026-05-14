package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		cache := NewCache(3)
		key := Key("key")
		key2 := Key("key2")
		key3 := Key("key3")
		key4 := Key("key4")
		value, value2, value3, value4 := "value", "value2", "value3", "value4"

		cache.Set(key, value)
		cache.Set(key2, value2)
		cache.Set(key3, value3)
		cache.Get(key3)
		cache.Get(key2)
		cache.Get(key)
		cache.Set(key4, value4)

		got3, ok3 := cache.Get(key3)
		require.False(t, ok3)
		require.Nil(t, got3)

		got4, ok4 := cache.Get(key4)
		require.True(t, ok4)
		require.Equal(t, value4, got4)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func Test_lruCache_Set(t *testing.T) {
	t.Run("new elem", func(t *testing.T) {
		capacity := 2
		cache := &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}
		key := Key("key")
		value := "value"

		wasInCache := cache.Set(key, value)

		require.False(t, wasInCache)
		require.Equal(t, 1, cache.queue.Len())
		require.Len(t, cache.items, 1)
		require.Contains(t, cache.items, key)
		require.Equal(t, key, cache.queue.Front().Key)
		require.Equal(t, value, cache.queue.Front().Value)
	})

	t.Run("existing elem", func(t *testing.T) {
		capacity := 2
		cache := &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}
		key := Key("key")
		oldValue := "value"
		newValue := "new value"
		cache.Set(key, oldValue)

		wasInCache := cache.Set(key, newValue)

		require.True(t, wasInCache)
		require.Equal(t, 1, cache.queue.Len())
		require.Len(t, cache.items, 1)
		require.Equal(t, key, cache.queue.Front().Key)
		require.Equal(t, newValue, cache.queue.Front().Value)
	})

	t.Run("existing elem moves to front", func(t *testing.T) {
		capacity := 3
		cache := &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}
		key := Key("key")
		key2 := Key("key2")
		key3 := Key("key3")
		cache.Set(key, "value")
		cache.Set(key2, "value2")
		cache.Set(key3, "value3")

		wasInCache := cache.Set(key, "new value")

		require.True(t, wasInCache)
		require.Equal(t, capacity, cache.queue.Len())
		require.Equal(t, key, cache.queue.Front().Key)
		require.Equal(t, key2, cache.queue.Back().Key)
		require.Equal(t, "new value", cache.queue.Front().Value)
	})

	t.Run("purge old elem", func(t *testing.T) {
		capacity := 2
		cache := &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}
		oldKey := Key("old")
		key := Key("key")
		newKey := Key("new")
		cache.Set(oldKey, "old value")
		cache.Set(key, "value")

		wasInCache := cache.Set(newKey, "new value")

		require.False(t, wasInCache)
		require.Equal(t, capacity, cache.queue.Len())
		require.Len(t, cache.items, capacity)
		require.NotContains(t, cache.items, oldKey)
		require.Contains(t, cache.items, key)
		require.Contains(t, cache.items, newKey)
		require.Equal(t, newKey, cache.queue.Front().Key)
		require.Equal(t, key, cache.queue.Back().Key)
	})
}

func Test_lruCache_deleteBack(t *testing.T) {
	t.Run("empty cash", func(t *testing.T) {
		capacity := 2
		cache := &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}

		require.NotPanics(t, func() { cache.deleteBack() })
	})

	t.Run("one elem, two cap", func(t *testing.T) {
		capacity := 2
		cache := &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}
		cache.Set("key", "value")
		cache.deleteBack()
		require.Equal(t, 1, cache.queue.Len())
	})

	t.Run("two elem, two cap", func(t *testing.T) {
		capacity := 2
		cache := &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}
		key := Key("key")
		key2 := Key("key2")
		cache.Set(key, "value")
		cache.Set(key2, "value2")
		cache.deleteBack()
		require.Equal(t, 1, cache.queue.Len())

		require.Contains(t, cache.items, key2)
	})
}

func Test_lruCache_Clear(t *testing.T) {
	t.Run("sucess", func(t *testing.T) {
		capacity := 5
		c := &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}
		wantLen := 2
		wantLenAfterClear := 0
		key := Key("key")
		key2 := Key("key2")
		c.Set(key, "value")
		c.Set(key2, "value2")
		require.Equal(t, wantLen, c.queue.Len())
		require.Len(t, c.items, wantLen)
		c.Clear()
		require.Equal(t, wantLenAfterClear, c.queue.Len())
		require.Len(t, c.items, wantLenAfterClear)
	})
}

func Test_lruCache_Get(t *testing.T) {
	t.Run("sucess", func(t *testing.T) {
		capacity := 5
		c := &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}
		wantValue := "value"
		wantValue2 := "value2"
		key := Key("key")
		key2 := Key("key2")
		c.Set(key, wantValue)
		c.Set(key2, wantValue2)
		value, ok := c.Get(key)
		value2, ok2 := c.Get(key2)
		require.Equal(t, wantValue, value)
		require.Equal(t, wantValue2, value2)
		require.True(t, ok)
		require.True(t, ok2)
	})
	t.Run("elem not in cashe", func(t *testing.T) {
		capacity := 5
		c := &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}
		key := Key("key")
		value, ok := c.Get(key)
		require.Nil(t, value)
		require.False(t, ok)
	})

	t.Run("last get elem become front", func(t *testing.T) {
		capacity := 5
		c := &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}
		wantValue := "value"
		wantValue2 := "value2"
		key := Key("key")
		key2 := Key("key2")
		c.Set(key, wantValue)
		c.Set(key2, wantValue2)
		_, _ = c.Get(key)
		value2, _ := c.Get(key2)
		require.Equal(t, value2, c.queue.Front().Value)
	})
}

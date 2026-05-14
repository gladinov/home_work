package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	value, value2, value3, value4 = "value", "value2", "value3", "value4"
	newValue                      = "new value"
	key                           = Key("key")
	key2                          = Key("key2")
	key3                          = Key("key3")
	key4                          = Key("key4")
	invalidCacheValue             = "invalid cache value"
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

func TestCacheMultithreading(_ *testing.T) {
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
		cache := &lruCache{
			capacity: 2,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 2),
		}

		wasInCache := cache.Set(key, value)

		require.False(t, wasInCache)
		require.Equal(t, 1, cache.queue.Len())
		require.Len(t, cache.items, 1)
		require.Contains(t, cache.items, key)
		require.Equal(t, key, cacheItemFromListItem(t, cache.queue.Front()).key)
		require.Equal(t, value, cacheItemFromListItem(t, cache.queue.Front()).value)
	})

	t.Run("existing elem", func(t *testing.T) {
		cache := &lruCache{
			capacity: 2,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 2),
		}

		cache.Set(key, value)

		wasInCache := cache.Set(key, newValue)

		require.True(t, wasInCache)
		require.Equal(t, 1, cache.queue.Len())
		require.Len(t, cache.items, 1)
		require.Equal(t, key, cacheItemFromListItem(t, cache.queue.Front()).key)
		require.Equal(t, newValue, cacheItemFromListItem(t, cache.queue.Front()).value)
	})

	t.Run("existing elem moves to front", func(t *testing.T) {
		cache := &lruCache{
			capacity: 3,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 3),
		}

		cache.Set(key, value)
		cache.Set(key2, value2)
		cache.Set(key3, value3)

		wasInCache := cache.Set(key, newValue)

		require.True(t, wasInCache)
		require.Equal(t, 3, cache.queue.Len())
		require.Equal(t, key, cacheItemFromListItem(t, cache.queue.Front()).key)
		require.Equal(t, key2, cacheItemFromListItem(t, cache.queue.Back()).key)
		require.Equal(t, newValue, cacheItemFromListItem(t, cache.queue.Front()).value)
	})

	t.Run("purge old elem", func(t *testing.T) {
		cache := &lruCache{
			capacity: 2,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 2),
		}
		oldKey := Key("old")
		newKey := Key("new")

		cache.Set(oldKey, "old value")
		cache.Set(key, value)

		wasInCache := cache.Set(newKey, newValue)

		require.False(t, wasInCache)
		require.Equal(t, 2, cache.queue.Len())
		require.Len(t, cache.items, 2)
		require.NotContains(t, cache.items, oldKey)
		require.Contains(t, cache.items, key)
		require.Contains(t, cache.items, newKey)
		require.Equal(t, newKey, cacheItemFromListItem(t, cache.queue.Front()).key)
		require.Equal(t, key, cacheItemFromListItem(t, cache.queue.Back()).key)
	})

	t.Run("capacity one evicts previous item", func(t *testing.T) {
		cache := NewCache(1)
		cache.Set("a", 1)
		cache.Set("b", 2)
		gotA, okA := cache.Get("a")
		require.False(t, okA)
		require.Nil(t, gotA)

		gotB, okB := cache.Get("b")
		require.True(t, okB)
		require.Equal(t, 2, gotB)
	})

	t.Run("zero capacity stores nothing", func(t *testing.T) {
		cache := NewCache(0)

		wasInCache := cache.Set(key, value)
		got, ok := cache.Get(key)

		require.False(t, wasInCache)
		require.False(t, ok)
		require.Nil(t, got)
	})

	t.Run("negative capacity stores nothing", func(t *testing.T) {
		cache := NewCache(-1)

		wasInCache := cache.Set(key, value)
		got, ok := cache.Get(key)

		require.False(t, wasInCache)
		require.False(t, ok)
		require.Nil(t, got)
	})
}

func Test_lruCache_deleteBack(t *testing.T) {
	t.Run("empty cash", func(t *testing.T) {
		cache := &lruCache{
			capacity: 2,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 2),
		}

		require.NotPanics(t, func() { cache.deleteBack() })
	})

	t.Run("one elem, two cap", func(t *testing.T) {
		cache := &lruCache{
			capacity: 2,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 2),
		}
		cache.Set(key, value)
		cache.deleteBack()
		require.Equal(t, 1, cache.queue.Len())
	})

	t.Run("two elem, two cap", func(t *testing.T) {
		cache := &lruCache{
			capacity: 2,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 2),
		}
		cache.Set(key, value)
		cache.Set(key2, value2)
		cache.deleteBack()
		require.Equal(t, 1, cache.queue.Len())

		require.Contains(t, cache.items, key2)
	})

	t.Run("panics on invalid cache item", func(t *testing.T) {
		cache := &lruCache{
			capacity: 1,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 1),
		}
		cache.queue.PushFront(invalidCacheValue)

		require.Panics(t, cache.deleteBack)
	})
}

func Test_lruCache_Clear(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := &lruCache{
			capacity: 5,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 5),
		}
		wantLen := 2
		wantLenAfterClear := 0
		c.Set(key, value)
		c.Set(key2, value2)
		require.Equal(t, wantLen, c.queue.Len())
		require.Len(t, c.items, wantLen)
		c.Clear()
		require.Equal(t, wantLenAfterClear, c.queue.Len())
		require.Len(t, c.items, wantLenAfterClear)
	})
}

func Test_lruCache_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := &lruCache{
			capacity: 5,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 5),
		}

		c.Set(key, value)
		c.Set(key2, value2)
		got, ok := c.Get(key)
		got2, ok2 := c.Get(key2)
		require.Equal(t, value, got)
		require.Equal(t, value2, got2)
		require.True(t, ok)
		require.True(t, ok2)
	})
	t.Run("elem not in cashe", func(t *testing.T) {
		c := &lruCache{
			capacity: 5,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 5),
		}
		value, ok := c.Get(key)
		require.Nil(t, value)
		require.False(t, ok)
	})

	t.Run("last get elem become front", func(t *testing.T) {
		c := &lruCache{
			capacity: 5,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 5),
		}

		c.Set(key, value)
		c.Set(key2, value2)
		_, _ = c.Get(key)
		got2, _ := c.Get(key2)
		require.Equal(t, got2, cacheItemFromListItem(t, c.queue.Front()).value)
	})

	t.Run("panics on invalid cache item", func(t *testing.T) {
		c := &lruCache{
			capacity: 1,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, 1),
		}
		item := c.queue.PushFront(invalidCacheValue)
		c.items[key] = item

		require.Panics(t, func() {
			_, _ = c.Get(key)
		})
	})
}

func cacheItemFromListItem(t *testing.T, item *ListItem) *cacheItem {
	t.Helper()
	require.NotNil(t, item)
	value, ok := item.Value.(*cacheItem)
	require.True(t, ok)
	return value
}

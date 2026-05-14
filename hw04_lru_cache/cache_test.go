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

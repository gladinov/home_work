package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront("10", 10) // [10]
		l.PushBack("20", 20)  // [10, 20]
		l.PushBack("30", 30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront("key", v)
			} else {
				l.PushBack("key", v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func Test_list_PushFront(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var l list
		key := Key("key")
		value := 10
		want := NewListItem(key, value, nil, nil)
		got := l.PushFront(key, value)
		require.Equal(t, want, got)
		require.Equal(t, want, l.Front())
		require.Equal(t, want, l.Back())
		require.Equal(t, 1, l.len)
	})

	t.Run("new front with one elem", func(t *testing.T) {
		var l list
		firstKey := Key("first")
		firstvalue := 0
		first := l.PushFront(firstKey, firstvalue)
		key := Key("key")
		value := 10
		want := NewListItem(key, value, first, nil)
		got := l.PushFront(key, value)
		require.Equal(t, want, got)
		require.Equal(t, want, l.Front())
		require.Equal(t, first, l.Front().Next)
		require.Equal(t, 2, l.len)
		require.NotNil(t, l.Back())
		require.Equal(t, first, l.Back())
	})

	t.Run("new front with several elem", func(t *testing.T) {
		var l list
		firstKey := Key("first")
		firstvalue := 1
		first := l.PushFront(firstKey, firstvalue)
		secondKey := Key("second")
		secondValue := 2
		second := l.PushFront(secondKey, secondValue)
		key := Key("key")
		value := 10
		want := NewListItem(key, value, second, nil)
		got := l.PushFront(key, value)
		require.Equal(t, want, got)
		require.Equal(t, want, l.Front())
		require.Equal(t, second, l.Front().Next)
		require.Equal(t, 3, l.Len())
		require.NotNil(t, l.Back())
		require.Equal(t, first, l.Back())
	})
}

func Test_list_PushBack(t *testing.T) {
	t.Run("zero list", func(t *testing.T) {
		var l list
		key := Key("key")
		value := 10
		want := NewListItem(key, value, nil, nil)
		got := l.PushBack(key, value)
		require.Equal(t, want, got)
		require.Equal(t, want, l.Front())
		require.Equal(t, want, l.Back())
		require.Equal(t, 1, l.len)
	})

	t.Run("one elem list", func(t *testing.T) {
		var l list
		firstKey := Key("first")
		firstvalue := 1
		first := l.PushFront(firstKey, firstvalue)
		key := Key("key")
		value := 10
		want := NewListItem(key, value, nil, first)
		got := l.PushBack(key, value)
		require.Equal(t, want, got)
		require.Equal(t, want, l.Back())
		require.Equal(t, 2, l.len)
		require.Equal(t, first, l.Front())
		require.Equal(t, want, l.Front().Next)
	})

	t.Run("more one elem list", func(t *testing.T) {
		var l list
		firstKey := Key("first")
		firstvalue := 1
		first := l.PushFront(firstKey, firstvalue)
		secondKey := Key("second")
		secondValue := 2
		second := l.PushFront(secondKey, secondValue)
		key := Key("key")
		value := 10
		want := NewListItem(key, value, nil, first)
		got := l.PushBack(key, value)
		require.Equal(t, want, got)
		require.Equal(t, want, l.Back())
		require.Equal(t, 3, l.len)
		require.Equal(t, second, l.Front())
		require.Equal(t, want, first.Next)
	})
}

func Test_list_Remove(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var l list

		require.NotPanics(t, func() { l.Remove(nil) })
	})

	t.Run("delete second elem", func(t *testing.T) {
		var l list
		wantLen := 2
		firstKey, secondKey, thirdKey := Key("first"), Key("second"), Key("third")
		firstValue, secondValue, thirdValue := 1, 2, 3
		first := l.PushFront(firstKey, firstValue)
		second := l.PushFront(secondKey, secondValue)
		third := l.PushFront(thirdKey, thirdValue) // [3,2,1]

		l.Remove(second)

		require.Equal(t, wantLen, l.Len())
		require.Equal(t, first, l.Front().Next)
		require.Equal(t, third, l.Back().Prev)
	})

	t.Run("delete back elem", func(t *testing.T) {
		var l list
		wantLen := 2
		firstKey, secondKey, thirdKey := Key("first"), Key("second"), Key("third")
		firstValue, secondValue, thirdValue := 1, 2, 3
		first := l.PushFront(firstKey, firstValue)
		second := l.PushFront(secondKey, secondValue)
		_ = l.PushFront(thirdKey, thirdValue) // [3,2,1]

		l.Remove(first)

		require.Equal(t, wantLen, l.Len())
		require.Equal(t, second, l.Back())
	})

	t.Run("delete front elem", func(t *testing.T) {
		var l list
		wantLen := 2
		firstKey, secondKey, thirdKey := Key("first"), Key("second"), Key("third")
		firstValue, secondValue, thirdValue := 1, 2, 3
		_ = l.PushFront(firstKey, firstValue)
		second := l.PushFront(secondKey, secondValue)
		third := l.PushFront(thirdKey, thirdValue) // [3,2,1]

		l.Remove(third)

		require.Equal(t, wantLen, l.Len())
		require.Equal(t, second, l.Front())
	})

	t.Run("delete single elem", func(t *testing.T) {
		var l list
		wantLen := 0
		firstKey := Key("first")
		firstValue := 1
		first := l.PushFront(firstKey, firstValue)

		l.Remove(first)

		require.Equal(t, wantLen, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})
}

func Test_list_MoveToFront(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var l list

		require.NotPanics(t, func() { l.MoveToFront(nil) })
		require.Equal(t, 0, l.Len())
	})

	t.Run("i already first", func(t *testing.T) {
		var l list
		firstKey, secondKey, thirdKey := Key("first"), Key("second"), Key("third")
		firstValue, secondValue, thirdValue := 1, 2, 3
		_ = l.PushFront(firstKey, firstValue)
		_ = l.PushFront(secondKey, secondValue)
		third := l.PushFront(thirdKey, thirdValue) // [3,2,1]
		require.NotPanics(t, func() { l.MoveToFront(third) })
		require.Equal(t, 3, l.Len())
		require.Equal(t, third, l.Front())
	})

	t.Run("_", func(t *testing.T) {
		var l list
		firstKey, secondKey, thirdKey := Key("first"), Key("second"), Key("third")
		firstValue, secondValue, thirdValue := 1, 2, 3
		first := l.PushBack(firstKey, firstValue)
		second := l.PushBack(secondKey, secondValue)
		third := l.PushBack(thirdKey, thirdValue)              // [1,2,3]
		require.NotPanics(t, func() { l.MoveToFront(second) }) // [2,1,3]
		require.Equal(t, 3, l.Len())
		require.Equal(t, second, l.Front())
		require.Equal(t, first.Prev, second)
		require.Equal(t, third.Prev, first)
	})
}

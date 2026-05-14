package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(key Key, v any) *ListItem
	PushBack(key Key, v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Key   Key
	Value any
	Next  *ListItem
	Prev  *ListItem
}

func NewListItem(key Key, value any, next, prev *ListItem) *ListItem {
	return &ListItem{
		Key:   key,
		Value: value,
		Next:  next,
		Prev:  prev,
	}
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(key Key, v any) *ListItem {
	newFront := NewListItem(key, v, l.front, nil)
	if l.Front() == nil {
		l.front = newFront
	} else {
		l.front.Prev = newFront
		if l.Len() == 1 {
			l.back = l.front
		}
	}
	l.front = newFront
	l.len++
	return l.front
}

func (l *list) PushBack(key Key, v any) *ListItem {
	if l.front == nil {
		l.front = NewListItem(key, v, nil, nil)
		l.len++
		return l.front
	}
	if l.Back() == nil {
		newBack := NewListItem(key, v, nil, l.front)
		l.front.Next = newBack
		l.back = newBack
	} else {
		newBack := NewListItem(key, v, nil, l.back)
		l.back.Next = newBack
		l.back = newBack
	}
	l.len++
	return l.back
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if i == l.Front() {
		l.front = i.Next
	}

	if i == l.Back() {
		l.back = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	if l.front == i {
		return
	}
	// Удалили элемент с текущей позиции
	l.Remove(i)

	// Задали первому элементу предудыший элемент равным i
	l.front.Prev = i
	// Предыдущий для i равен nil после перемещения его впредед
	i.Prev = nil
	// Задали новому первому элементу next
	i.Next = l.front
	// Сделали i первым элементом list
	l.front = i

	// Вернули счетчик к исзодному значению
	l.len++
}

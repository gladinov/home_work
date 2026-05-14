package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

func NewListItem(value any, next, prev *ListItem) *ListItem {
	return &ListItem{
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

func (l *list) PushFront(v any) *ListItem {
	newFront := NewListItem(v, l.front, nil)
	if l.Front() == nil {
		l.back = newFront
	} else {
		l.front.Prev = newFront
	}
	l.front = newFront
	l.len++
	return l.front
}

func (l *list) PushBack(v any) *ListItem {
	if l.front == nil {
		newFront := NewListItem(v, nil, nil)
		l.front = newFront
		l.back = newFront
	} else {
		newBack := NewListItem(v, nil, l.back)
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
	l.Remove(i)

	l.front.Prev = i
	i.Prev = nil
	i.Next = l.front
	l.front = i

	l.len++
}

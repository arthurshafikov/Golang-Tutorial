package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front, back *ListItem
	len         int
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	front := l.Front()

	newItem := &ListItem{
		v, front, nil,
	}

	if front != nil {
		front.Prev = newItem
	} else {
		l.back = newItem
	}
	l.len++
	l.front = newItem

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	back := l.Back()

	newItem := &ListItem{
		v, nil, back,
	}

	if back != nil {
		back.Next = newItem
	} else {
		l.front = newItem
	}
	l.len++
	l.back = newItem

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}

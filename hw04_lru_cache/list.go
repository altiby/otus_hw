package hw04lrucache

import "sync"

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
	m     sync.Mutex
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	l.m.Lock()
	defer l.m.Unlock()
	return l.len
}

func (l *list) Front() *ListItem {
	l.m.Lock()
	defer l.m.Unlock()
	return l.front
}

func (l *list) Back() *ListItem {
	l.m.Lock()
	defer l.m.Unlock()
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.m.Lock()
	defer l.m.Unlock()
	return l.pushFront(v)
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.m.Lock()
	defer l.m.Unlock()
	return l.pushBack(v)
}

func (l *list) Remove(i *ListItem) {
	l.m.Lock()
	defer l.m.Unlock()
	l.remove(i)
}

func (l *list) MoveToFront(i *ListItem) {
	l.m.Lock()
	defer l.m.Unlock()
	l.moveToFront(i)
}

func NewList() List {
	return new(list)
}

func (l *list) pushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  l.front,
	}
	if l.front != nil {
		l.front.Prev = newItem
	}
	l.front = newItem
	if l.back == nil {
		l.back = newItem
	}
	l.len++
	return newItem
}

func (l *list) pushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Prev:  l.back,
	}
	if l.back != nil {
		l.back.Next = newItem
	}
	l.back = newItem
	if l.front == nil {
		l.front = newItem
	}
	l.len++
	return newItem
}

func (l *list) remove(i *ListItem) {
	if l.len == 0 {
		return
	}

	if l.len == 1 {
		l.front = nil
		l.back = nil
		return
	}

	l.len--

	if l.front == i {
		l.front = l.front.Next
		l.front.Prev = nil
		return
	}

	if l.back == i {
		l.back = l.back.Prev
		l.back.Next = nil
		return
	}

	i.Prev.Next, i.Next.Prev = i.Next, i.Prev
}

func (l *list) moveToFront(i *ListItem) {
	if l.len < 1 {
		return
	}
	if i == l.front {
		return
	}

	l.remove(i)
	i.Prev = nil
	i.Next = nil
	l.len++
	i.Next = l.front
	i.Prev = nil
	l.front.Prev = i
	l.front = i
}

package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	m        sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheData struct {
	CacheKey   Key
	CacheValue interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.m.Lock()
	defer l.m.Unlock()
	elem, ok := l.items[key]
	if ok {
		elem.Value = &cacheData{
			CacheKey:   key,
			CacheValue: value,
		}
		l.queue.MoveToFront(elem)
	} else {
		elem = l.queue.PushFront(&cacheData{
			CacheKey:   key,
			CacheValue: value,
		})
		l.items[key] = elem
		if l.queue.Len() > l.capacity {
			back := l.queue.Back()
			l.queue.Remove(back)
			delete(l.items, back.Value.(*cacheData).CacheKey)
		}
	}
	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.m.Lock()
	defer l.m.Unlock()
	elem, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(elem)
		return elem.Value.(*cacheData).CacheValue, ok
	}
	return nil, ok
}

func (l *lruCache) Clear() {
	l.m.Lock()
	defer l.m.Unlock()
	l.items = make(map[Key]*ListItem, l.capacity)
	l.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

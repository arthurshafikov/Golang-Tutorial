package hw04lrucache

import "sync"

type Key string

type Items map[Key]*ListItem

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    Items
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	newItem := &cacheItem{key, value}
	if _, ok := l.items[key]; ok { // item exists
		l.items[key].Value = newItem
		l.queue.MoveToFront(l.items[key])

		return true
	}

	listItem := l.queue.PushFront(newItem)
	l.items[key] = listItem
	if l.queue.Len() > l.capacity {
		lastElem := l.queue.Back()
		l.queue.Remove(lastElem)
		delete(l.items, lastElem.Value.(*cacheItem).key)
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	elem, ok := l.items[key]
	if !ok {
		return nil, false
	}
	l.queue.MoveToFront(elem)

	return elem.Value.(*cacheItem).value, true
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = NewList()
	l.items = Items{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(Items, capacity),
	}
}

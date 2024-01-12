package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	isFilled() bool
	Clear()
}

type lruCache struct {
	mutex    sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type Item struct {
	Key   Key
	Value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		item.Value.(*Item).Value = value
		return true
	}

	item := &Item{
		Key:   key,
		Value: value,
	}

	if l.queue.Len() == l.capacity {
		l.purge()
	}

	l.queue.PushFront(item)
	l.items[key] = l.queue.Front()

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		return item.Value.(*Item).Value, true
	}

	return nil, false
}

func (l *lruCache) isFilled() bool {
	return len(l.items) == l.capacity
}

func (l *lruCache) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func (l *lruCache) purge() {
	item := l.queue.Back()
	delete(l.items, item.Value.(*Item).Key)
	l.queue.Remove(item)
}

package cache

import (
	"sync"

	"github.com/Georgiiagon/image-previewer/internal/app"
	"github.com/Georgiiagon/image-previewer/internal/cache/lrucache"
)

type Cache interface {
	Set(key app.Key, value interface{}) bool
	Get(key app.Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    lrucache.List
	items    map[app.Key]*lrucache.ListItem
}

type cacheItem struct {
	key   app.Key
	value interface{}
}

func (lru *lruCache) Set(key app.Key, value interface{}) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	item, ok := lru.items[key]
	cItem := cacheItem{key: key, value: value}

	if ok {
		item.Value = cItem
		lru.queue.MoveToFront(item)
	} else {
		lru.queue.PushFront(cItem)
	}

	if lru.queue.Len() > lru.capacity {
		lru.deleteLast()
	}

	lru.items[key] = lru.queue.Front()

	return ok
}

func (lru *lruCache) Get(key app.Key) (interface{}, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	item, ok := lru.items[key]
	if !ok {
		return nil, false
	}

	lru.queue.MoveToFront(item)
	cItem := item.Value.(cacheItem)

	return cItem.value, true
}

func (lru *lruCache) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	lru.items = make(map[app.Key]*lrucache.ListItem, lru.capacity)
	lru.queue = lrucache.NewList()
}

func (lru *lruCache) deleteLast() {
	lastItem := lru.queue.Back()
	lru.queue.Remove(lastItem)
	cItem := lastItem.Value.(cacheItem)
	delete(lru.items, cItem.key)
}

func New(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    lrucache.NewList(),
		items:    make(map[app.Key]*lrucache.ListItem, capacity),
	}
}

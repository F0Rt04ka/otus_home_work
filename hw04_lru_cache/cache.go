package hw04lrucache

import "sync"

type Key string

type CacheItem struct {
	Key   Key
	Value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mutex    sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.items[key]
	if ok {
		item.Value.(*CacheItem).Value = value
		c.queue.MoveToFront(item)
		return true
	}

	if c.queue.Len() >= c.capacity {
		backItem := c.queue.Back()

		delete(c.items, backItem.Value.(*CacheItem).Key)
		c.queue.Remove(backItem)
	}

	c.items[key] = c.queue.PushFront(&CacheItem{
		Key:   key,
		Value: value,
	})

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(c.items[key])

	return item.Value.(*CacheItem).Value, true
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
